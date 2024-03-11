package apiconfig

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/kafka-go"
)

// getKafkaTopics returns a list of topics from Kafka
func GetKafkaTopics() ([]string, error) {
	conn, err := kafka.Dial("tcp", "kafka:9092") // Adjust according to your Kafka setup
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return nil, err
	}

	topicMap := make(map[string]struct{})
	for _, p := range partitions {
		topicMap[p.Topic] = struct{}{}
	}

	var topics []string
	for topic := range topicMap {
		topics = append(topics, topic)
	}

	return topics, nil
}

var jwtKey = []byte("swg3s+hZkEz/Vh7fXtJRCgRAJBzT2ttyHcIgwD14b3k=")

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

	c.Next()   // call Next and don't return its value
	return nil // return nil since we don't have any errors to report
}
