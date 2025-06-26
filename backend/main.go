package main

import (
	"log"
	"net/http"

	"github.com/logoes0/peeriodic.git/handlers"
	"github.com/logoes0/peeriodic.git/routers"
)

func main() {
	routers.SetupRoutes()

	go handlers.HandleMessages()

	port := ":5000"
	log.Printf("Server listening on port %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
