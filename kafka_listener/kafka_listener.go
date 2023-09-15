package main

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func handleMessage(msg *kafka.Message) {
	fmt.Printf("A new message was received: %s\n", string(msg.Value))
}

func main() {
	fmt.Println("Hello World!")

	// Creating a Kafka consumer instance
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Error when creating a consumer: %s\n", err.Error())
		return
	}
	// Subscribing to a topic
	err = consumer.SubscribeTopics([]string{"my-topic"}, nil)
	if err != nil {
		fmt.Printf("Error when subscribing to a topic: %s\n", err.Error())
		return
	}
	// Message receiving cycle
	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			fmt.Printf("Error when reading a message: %s\n", err.Error())
			continue
		}
		if msg != nil {
			// Calling a handler function for each received message
			handleMessage(msg)
		} else {
			fmt.Printf("No message received\n")
		}
		// Closing the connection with the Kafka server
		consumer.Close()
	}
}
