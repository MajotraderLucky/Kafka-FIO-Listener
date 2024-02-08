package main

import (
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/MajotraderLucky/Kafka-FIO-Listener/checker"
	"github.com/MajotraderLucky/Utils/logger"
)

func main() {
	// Wait for other services to be ready
	time.Sleep(time.Second * 10)

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

	log.Println("Hello, logger!")

	// Connect to the database and check it
	db, err := checker.ConnectAndCheckDB()
	if err != nil {
		log.Fatal(err)
	}

	// Close the database connection
	defer db.Close()

	err = checker.CheckZookeeper()
	if err != nil {
		log.Fatalf("Error checking Zookeeper: %s", err)
	}
	log.Println("Zookeeper is working correctly")

	// Check Kafka service
	err = checker.CheckKafka()
	if err != nil {
		log.Fatalf("Error checking Kafka: %s", err)
	}
	log.Println("Kafka is working correctly")

	logger.CleanLogCountLines(50)
}
