package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/logoes0/peeriodic.git/handlers"
	"github.com/logoes0/peeriodic.git/routers"
)

func main() {
	db, err := sql.Open("postgres", "postgres://username:password@localhost:5432/collab_editor?sslmode=disable")
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
