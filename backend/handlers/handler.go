package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/logoes0/peeriodic.git/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow CORS
}

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan models.Message)
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

	// Get or create the room
	roomsMu.Lock()
	room, exists := rooms[roomID]
	if !exists {
		room = &models.Room{
			Clients:  make(map[*websocket.Conn]bool),
			Document: "", // initialize blank
		}
		rooms[roomID] = room
	}
	roomsMu.Unlock()

	// Add client to the room
	room.Mu.Lock()
	room.Clients[ws] = true

	// Send current document content to the new client
	ws.WriteJSON(models.Message{Type: "init", Data: room.Document})
	room.Mu.Unlock()

	log.Println("New client connected to room:", roomID)

	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading JSON: %v", err)
			room.Mu.Lock()
			delete(room.Clients, ws)
			room.Mu.Unlock()
			break
		}

		if msg.Type == "update" {
			// Store new content
			room.Mu.Lock()
			room.Document = msg.Data

			// Broadcast to all clients in the room
			for client := range room.Clients {
				if client != ws { // optional: prevent echo
					client.WriteJSON(msg)
				}
			}
			room.Mu.Unlock()
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

func HandleRooms(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodGet:
		uid := r.URL.Query().Get("uid")
		if uid == "" {
			http.Error(w, "Missing uid", http.StatusBadRequest)
			return
		}

		rows, err := db.Query("SELECT id, title, created_at FROM rooms WHERE user_uid=$1", uid)
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var rooms []map[string]interface{}
		for rows.Next() {
			var id, title string
			var createdAt string
			rows.Scan(&id, &title, &createdAt)
			rooms = append(rooms, map[string]interface{}{
				"id":         id,
				"title":      title,
				"created_at": createdAt,
			})
		}

		json.NewEncoder(w).Encode(rooms)

	case http.MethodPost:
		var payload struct {
			Title string `json:"title"`
			UID   string `json:"uid"`
		}
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil || payload.UID == "" {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		roomID := uuid.New()
		_, err = db.Exec("INSERT INTO rooms (id, title, user_uid) VALUES ($1, $2, $3)", roomID, payload.Title, payload.UID)
		if err != nil {
			http.Error(w, "Failed to create room", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"id":    roomID.String(),
			"title": payload.Title,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleRoomByID(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/api/rooms/"):]
	if id == "" {
		http.Error(w, "Missing room ID", http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT id, title, content FROM rooms WHERE id=$1", id)
	var roomID, title, content string
	err := row.Scan(&roomID, &title, &content)
	if err == sql.ErrNoRows {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"id":      roomID,
		"title":   title,
		"content": content,
	})
}
