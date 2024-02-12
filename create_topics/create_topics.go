package main

import (
	"context"
	"log"
	"time"

	"github.com/MajotraderLucky/Utils/logger"
	"github.com/segmentio/kafka-go"
)

func main() {
	logger := logger.Logger{}
	err := logger.CreateLogsDir()
	if err != nil {
		log.Fatal(err)
	}
	err = logger.OpenLogFile()
	if err != nil {
		log.Fatal(err)
	}
	logger.SetLogger()
	logger.LogLine()
	log.Println("Create_topics is starting...")
	logger.LogLine()

	// Kafka broker
	brokerAddress := "kafka:9092"

	// Creating a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Establishing a connection to the Kafka broker using the context
	conn, err := kafka.DialContext(ctx, "tcp", brokerAddress)
	if err != nil {
		log.Fatalf("failed to dial Kafka broker: %v", err)
	}
	defer conn.Close()

	// Define topics to be created
	topics := []kafka.TopicConfig{
		{
			Topic:             "example-topic",
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
		// Add more topics here
	}

	// Create topics
	err = conn.CreateTopics(topics...)
	if err != nil {
		log.Fatalf("failed to create topics: %v", err)
	}
	log.Println("Topics created successfully")
	logger.LogLine()
}
