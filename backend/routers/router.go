package routers

import (
	"net/http"

	"github.com/logoes0/peeriodic.git/handlers"
)

func SetupRoutes() {
	http.HandleFunc("/ws", handlers.HandleConnections)
}
