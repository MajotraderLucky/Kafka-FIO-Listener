package main

import (
	"fmt"
	"log"
	"time"

	"database/sql"

	_ "github.com/lib/pq"

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

	// Connect to the database
	db, err := sql.Open("postgres", "host=db port=5432 user=postgres password=mysecretpassword dbname=mydatabase sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// Check if the database exists
	_, err = db.Exec("SELECT 1")
	if err != nil {
		log.Fatal("Database does not exist: ", err)
	} else {
		log.Println("Database exists")
	}

	// Check if the table exists
	_, err = db.Exec("SELECT 1 FROM fio_data LIMIT 1")
	if err != nil {
		log.Fatal("Table does not exist: ", err)
	} else {
		log.Println("Table exists")
	}

	// Close the database connection
	defer db.Close()
}
