package checker

import (
	"database/sql"
	"fmt"
	"log"
)

func Hello() {
	log.Println("Hello, checker!")
}

func ConnectAndCheckDB() (*sql.DB, error) {
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
