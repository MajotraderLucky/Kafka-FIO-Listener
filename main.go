package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/MajotraderLucky/Kafka-FIO-Listener/checker"
	"github.com/MajotraderLucky/Utils/logger"
	zk "github.com/go-zookeeper/zk"
)

// CheckZookeeper checks the connection and operations with Zookeeper
func CheckZookeeper() error {
	conn, _, err := zk.Connect([]string{"zookeeper:2181"}, 10*time.Second)
	if err != nil {
		log.Println("Failed to connect to Zookeeper:", err)
		return fmt.Errorf("failed to connect to Zookeeper: %s", err)
	}
	defer conn.Close()

	_, err = conn.Create("/test", []byte("test"), 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Println("Failed to create test node:", err)
		return fmt.Errorf("failed to create test node: %s", err)
	}

	defer func() {
		conn.Delete("/test", -1)
	}()

	data, stat, err := conn.Get("/test")
	if err != nil {
		log.Println("Failed to get data from test node:", err)
		return fmt.Errorf("failed to get data from test node: %s", err)
	}

	if string(data) != "test" {
		log.Println("znode content does not match expected:", data)
		return fmt.Errorf("znode content does not match expected: %s", data)
	}

	if stat.Cversion != 0 {
		log.Println("node version is not zero:", stat.Cversion)
		return fmt.Errorf("node version is not zero: %d", stat.Cversion)
	}

	eventChan := make(chan zk.Event)
	go func() {
		eventChan <- zk.Event{Type: zk.EventSession}
	}()

	select {
	case event := <-eventChan:
		if event.Type != zk.EventSession {
			log.Println("expected EventSession type event, got:", event.Type)
			return fmt.Errorf("expected EventSession type event, got: %v", event.Type)
		}
	case <-time.After(time.Second * 10):
		log.Println("timeout while waiting for an event")
		return fmt.Errorf("timeout while waiting for an event")
	}

	return nil
}

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

	err = CheckZookeeper()
	if err != nil {
		log.Fatalf("Error checking Zookeeper: %s", err)
	}
	log.Println("Zookeeper is working correctly")

	logger.CleanLogCountLines(50)
}
