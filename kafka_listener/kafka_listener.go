package main

import (
	"log"
	"time"

	"github.com/MajotraderLucky/Utils/logger"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func handleMessage(msg *kafka.Message) {
	log.Printf("A new message was received: %s\n", string(msg.Value))
}

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

	// Add a delay of 5 seconds
	time.Sleep(10 * time.Second)

	logger.LogLine()
	log.Println("Kafka listener is Starting...")
	logger.LogLine()

	// Creating a Kafka consumer instance
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	})

	// Check if connection to Kafka was successful
	if err != nil {
		log.Printf("Error when creating a consumer: %s\n", err.Error())
		return
	} else {
		log.Println("Successfully connected to Kafka from Kafka listener service")
	}

	// Subscribing to a topic
	err = consumer.SubscribeTopics([]string{"my-topic"}, nil)
	if err != nil {
		log.Printf("Error when subscribing to a topic: %s\n", err.Error())
		return
	}

	// Message receiving cycle
	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Error when reading a message: %s\n", err.Error())
			continue
		}
		if msg != nil {
			// Calling a handler function for each received message
			handleMessage(msg)
		} else {
			log.Printf("No message received\n")
		}
		// Closing the connection with the Kafka server
		consumer.Close()
	}
}
