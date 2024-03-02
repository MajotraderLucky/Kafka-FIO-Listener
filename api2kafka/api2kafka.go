package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/MajotraderLucky/Utils/logger"

	"github.com/gofiber/fiber/v2"
)

// Structure to store the response from the API
type ApiResponse struct {
	Data string `json:"data"`
}

func fetchDataFromAPI(url string) (string, error) {
	// Send GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// Parse JSON
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return "", err
	}
	return apiResponse.Data, nil
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
	logger.LogLine()
	log.Println("API2Kafka is starting...")
	logger.LogLine()

	app := fiber.New()

	app.Get("data", func(c *fiber.Ctx) error {
		// URL of the open API
		url := "http://127.0.0.1:8086/data"

		data, err := fetchDataFromAPI(url)
		if err != nil {
			log.Println(err)
			return c.Status(500).SendString("Internal server error")
		}
		// Display the result
		return c.SendString(data)
	})
	// Start server
	app.Listen(":3000")
}
