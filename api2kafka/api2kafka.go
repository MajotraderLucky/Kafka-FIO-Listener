package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/IBM/sarama"
	"github.com/MajotraderLucky/Kafka-FIO-Listener/apiconfig"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// KafkaInfo is a structure for storing information about Kafka topics
type KafkaInfo struct {
	Topics []string `json:"topics"`
}

type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("swg3s+hZkEz/Vh7fXtJRCgRAJBzT2ttyHcIgwD14b3k=")

// TopicsMessagesCount is a structure for storing information about the number of messages in Kafka topics
type TopicsMessagesCount struct {
	Topic               string `json:"topic"`
	TopicsMessagesCount int64  `json:"messagesCount"`
}

func main() {
	app := fiber.New()

	// Protecting the endpoint with the help of the isAuthorized middleware
	app.Get("/kafka/topics", IsAuthorized, func(c *fiber.Ctx) error {
		topics, err := apiconfig.GetKafkaTopics()
		if err != nil {
			log.Println("Failed to get Kafka topics:", err)
			return c.Status(http.StatusInternalServerError).SendString("Failed to get Kafka topics")
		}

		return c.JSON(KafkaInfo{Topics: topics})
	})

	app.Get("/kafka/topics/messages", IsAuthorized, func(c *fiber.Ctx) error {
		topics, err := apiconfig.GetKafkaTopics()
		if err != nil {
			log.Println("Failed to get Kafka topics:", err)
			return c.Status(http.StatusInternalServerError).SendString("Failed to get Kafka topics")
		}

		messagesCount, err := GetMessagesCount(topics)
		if err != nil {
			log.Println("Failed to get messages count:", err)
			return c.Status(http.StatusInternalServerError).SendString("Failed to get messages count")
		}

		return c.JSON(messagesCount)
	})

	err := app.Listen(":8086")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func IsAuthorized(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	return c.Next()
}

func GetMessagesCount(topics []string) ([]TopicsMessagesCount, error) {
	var topicsMessagesCount []TopicsMessagesCount

	// Initialize a new Sarama client configuration.
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0 // Specify the Kafka version here

	// Create a new Sarama client
	client, err := sarama.NewClient([]string{"kafka:9092"}, config)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	for _, topic := range topics {
		partitions, err := client.Partitions(topic)
		if err != nil {
			return nil, err
		}
		var totalMessages int64 = 0
		for _, partition := range partitions {
			oldestOffset, err := client.GetOffset(topic, partition, sarama.OffsetOldest)
			if err != nil {
				log.Println("Error getting partitions or offsets for topic", topic, ":", err)
				return nil, err
			}
			newestOffset, err := client.GetOffset(topic, partition, sarama.OffsetNewest)
			if err != nil {
				log.Println("Error getting partitions or offsets for topic", topic, ":", err)
				return nil, err
			}
			totalMessages += newestOffset - oldestOffset
		}
		topicsMessagesCount = append(topicsMessagesCount, TopicsMessagesCount{Topic: topic, TopicsMessagesCount: totalMessages})
	}

	return topicsMessagesCount, nil
}
