package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/MajotraderLucky/Kafka-FIO-Listener/apiconfig"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// KafkaInfo is a structure for storing information about Kafka topics
type KafkaInfo struct {
	Topics []string `json:"topics"`
}

type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("swg3s+hZkEz/Vh7fXtJRCgRAJBzT2ttyHcIgwD14b3k=")

func main() {
	app := fiber.New()

	// Protecting the endpoint with the help of the isAuthorized middleware
	app.Get("/kafka/topics", IsAuthorized, func(c *fiber.Ctx) error {
		topics, err := apiconfig.GetKafkaTopics()
		if err != nil {
			log.Println("Failed to get Kafka topics:", err)
			return c.Status(http.StatusInternalServerError).SendString("Failed to get Kafka topics")
		}

		return c.JSON(KafkaInfo{Topics: topics})
	})

	err := app.Listen(":8086")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func IsAuthorized(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	return c.Next()
}
