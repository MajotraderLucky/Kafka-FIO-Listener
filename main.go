package main

import (
	"fmt"
	"log"
	"time"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/MajotraderLucky/Utils/logger"
)

func connectAndCheckDB() (*sql.DB, error) {
	// Подключение к базе данных
	db, err := sql.Open("postgres", "host=db port=5432 user=postgres password=mysecretpassword dbname=mydatabase sslmode=disable")
	if err != nil {
		log.Println("ошибка при подключении к базе данных: ", err)
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %w", err)
	}

	// Проверка наличия базы данных
	_, err = db.Exec("SELECT 1")
	if err != nil {
		log.Println("база данных не существует: ", err)
		return nil, fmt.Errorf("база данных не существует: %w", err)
	} else {
		log.Println("база данных существует")
	}

	// Проверка наличия таблицы
	_, err = db.Exec("SELECT 1 FROM fio_data LIMIT 1")
	if err != nil {
		log.Println("таблица не существует: ", err)
		return nil, fmt.Errorf("таблица не существует: %w", err)
	} else {
		log.Println("таблица существует")
	}

	return db, nil
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
	db, err := connectAndCheckDB()
	if err != nil {
		log.Fatal(err)
	}

	// Close the database connection
	defer db.Close()

	logger.CleanLogCountLines(50)
}
