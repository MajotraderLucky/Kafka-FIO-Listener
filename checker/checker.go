package checker

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	zk "github.com/go-zookeeper/zk"
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
