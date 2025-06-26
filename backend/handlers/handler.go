package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/logoes0/peeriodic.git/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow CORS
}

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan models.Message)
	document  string
	mu        sync.Mutex
	rooms     = make(map[string]*models.Room)
	roomsMu   sync.Mutex
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		http.Error(w, "Missing room ID", http.StatusBadRequest)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket Upgrade error: %v", err)
		return
	}
	defer ws.Close()

	clients[ws] = true
	log.Println("New client connected")

	ws.WriteJSON(models.Message{Type: "init", Data: document})

	roomsMu.Lock()
	room, exists := rooms[roomID]
	if !exists {
		room = &models.Room{Clients: make(map[*websocket.Conn]bool)}
		rooms[roomID] = room
	}
	roomsMu.Unlock()

	room.Mu.Lock()
	room.Clients[ws] = true
	ws.WriteJSON(models.Message{Type: "init", Data: room.Document})
	room.Mu.Unlock()

	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading JSON: %v", err)
			delete(clients, ws)
			break
		}

		if msg.Type == "update" {
			mu.Lock()
			document = msg.Data
			mu.Unlock()
			broadcast <- msg
		}
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error broadcasting to client: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
