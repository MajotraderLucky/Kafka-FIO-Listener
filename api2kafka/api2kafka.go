package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"
)

// KafkaInfo is a structure for storing information about Kafka topics
type KafkaInfo struct {
	Topics []string `json:"topics"`
}

func main() {
	app := fiber.New()

	// Endpoint for fetching the list of topics
	app.Get("/kafka/topics", func(c *fiber.Ctx) error {
		topics, err := getKafkaTopics()
		if err != nil {
			log.Println("Failed to get Kafka topics:", err)
			return c.Status(http.StatusInternalServerError).SendString("Failed to get Kafka topics")
		}

		return c.JSON(KafkaInfo{Topics: topics})
	})

	err := app.Listen(":8086")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getKafkaTopics returns a list of topics from Kafka
func getKafkaTopics() ([]string, error) {
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
