package main

import (
	"log"
	"net/http"

	"github.com/MajotraderLucky/Kafka-FIO-Listener/apiconfig"
	"github.com/gofiber/fiber/v2"
)

// KafkaInfo is a structure for storing information about Kafka topics
type KafkaInfo struct {
	Topics []string `json:"topics"`
}

func main() {
	app := fiber.New()

	// Endpoint for fetching the list of topics
	app.Get("/kafka/topics", func(c *fiber.Ctx) error {
		topics, err := apiconfig.GetKafkaTopics()
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
