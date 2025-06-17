package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/logoes0/peeriodic.git/routes"
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

	// Setup all routes
	routes.SetupRoutes(appFiber, authClient)

	localURL := "http://localhost:" + port
	fmt.Println("🚀 Server running on:", localURL)
	log.Fatal(appFiber.Listen(":" + port))
}
