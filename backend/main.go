package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/logoes0/peeriodic.git/handlers"
	"github.com/logoes0/peeriodic.git/routers"
)

func main() {
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf(
		"postgres://%s@localhost:%s/%s?sslmode=disable",
		dbUser,
		dbPort,
		dbName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	routers.SetupRoutes(db)

	go handlers.HandleMessages()

	port := ":5000"
	log.Printf("Server listening on port %s", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
