package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	"github.com/logoes0/peeriodic.git/server"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("BE_PORT")

	authClient, err := server.InitFirebaseAuth()
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	appFiber := fiber.New()

	// basic api endpoint to check working
	appFiber.Get("/start", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "initial run successful",
		})
	})

	// middleware to upgrade only WebSocket requests
	appFiber.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// websocket endpoint
	appFiber.Get("/ws", websocket.New(server.HandleWebSocket(authClient)))

	localURL := "http://localhost:" + port
	fmt.Println("🚀 Server running on:", localURL)
	log.Fatal(appFiber.Listen(":" + port))
}
