package main

import (
	"fmt"
	"log"
	"time"

	"github.com/MajotraderLucky/Utils/logger"
)

func main() {
	// Wait for other services to be ready
	time.Sleep(time.Second * 10)

	fmt.Println("Hello, world!")

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
}
