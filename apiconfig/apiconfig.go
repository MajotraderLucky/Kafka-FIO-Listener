package apiconfig

import (
	"log"

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

func Hello() {
	log.Println("Hello, apiconfig!")
}
