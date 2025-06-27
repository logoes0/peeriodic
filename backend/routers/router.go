package routers

import (
	"database/sql"
	"net/http"

	"github.com/logoes0/peeriodic.git/handlers"
)

func SetupRoutes(db *sql.DB) {
	http.HandleFunc("/ws", handlers.HandleConnections)
	http.HandleFunc("/api/rooms", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRooms(w, r, db)
	}) // GET (list), POST (create)
	http.HandleFunc("/api/rooms/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			handlers.DeleteRoom(w, r, db)
			return
		}
		handlers.HandleRoomByID(w, r, db)
	}) // GET (by ID)
}
