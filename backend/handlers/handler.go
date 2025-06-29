package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/logoes0/peeriodic.git/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow CORS
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan models.Message)
	rooms     = make(map[string]*models.Room)
	roomsMu   sync.Mutex
)

func HandleConnections(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// 1. Validate room ID
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		http.Error(w, "Missing room ID", http.StatusBadRequest)
		return
	}

	// 2. Verify room exists or create it
	var content string
	err := db.QueryRow(`
        INSERT INTO rooms (id, content) 
        VALUES ($1, '') 
        ON CONFLICT (id) DO UPDATE SET id = EXCLUDED.id 
        RETURNING content
    `, roomID).Scan(&content)

	if err != nil {
		log.Printf("DB room initialization error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// 3. Upgrade to WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer func() {
		ws.Close()
		log.Printf("Client disconnected from room %s", roomID)
	}()

	// 4. Initialize room state
	roomsMu.Lock()
	room, exists := rooms[roomID]
	if !exists {
		room = &models.Room{
			Clients:  make(map[*websocket.Conn]bool),
			Document: content,
		}
		rooms[roomID] = room
	}
	roomsMu.Unlock()

	// 5. Add client and send current document
	room.Mu.Lock()
	room.Clients[ws] = true
	err = ws.WriteJSON(models.Message{
		Type: "init",
		Data: room.Document,
	})
	room.Mu.Unlock()

	if err != nil {
		log.Printf("Failed to send init message: %v", err)
		return
	}

	log.Printf("New client connected to room %s (%d clients)", roomID, len(room.Clients))

	// 6. Message handling loop
	for {
		var msg models.Message
		if err := ws.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Client disconnected unexpectedly: %v", err)
			} else {
				log.Printf("Message read error: %v", err)
			}
			break
		}

		switch msg.Type {
		case "update":
			room.Mu.Lock()

			// Update local document
			room.Document = msg.Data

			// Persist to database with transaction
			go func(content string) {
				tx, err := db.Begin()
				if err != nil {
					log.Printf("Transaction begin error: %v", err)
					return
				}
				defer tx.Rollback()

				if _, err := tx.Exec(`
                    UPDATE rooms 
                    SET content = $1, updated_at = NOW() 
                    WHERE id = $2
                `, content, roomID); err != nil {
					log.Printf("DB update error: %v", err)
					return
				}

				if err := tx.Commit(); err != nil {
					log.Printf("Transaction commit error: %v", err)
				}
			}(msg.Data)

			// Broadcast to other clients
			for client := range room.Clients {
				if client != ws {
					if err := client.WriteJSON(msg); err != nil {
						log.Printf("Broadcast error: %v", err)
						client.Close()
						delete(room.Clients, client)
					}
				}
			}
			room.Mu.Unlock()

		default:
			log.Printf("Unknown message type: %s", msg.Type)
		}
	}

	// 7. Clean up on disconnect
	room.Mu.Lock()
	delete(room.Clients, ws)
	remaining := len(room.Clients)
	room.Mu.Unlock()

	log.Printf("Client disconnected from room %s (%d clients remaining)", roomID, remaining)
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

func DeleteRoom(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	roomID := strings.TrimPrefix(r.URL.Path, "/api/rooms/")
	if roomID == "" {
		http.Error(w, "Missing room ID", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM rooms WHERE id = $1", roomID)
	if err != nil {
		http.Error(w, "Failed to delete room", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func HandleSave(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	enableCORS(&w)
	// Handle OPTIONS for preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.Println("Save endpoint hit") // Debug log

	// 1. Validate request method
	if r.Method != http.MethodPost {
		log.Println("Wrong method used") // Debug log
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Get room ID
	roomID := r.URL.Query().Get("room")
	log.Printf("Room ID from query: %s", roomID) // Debug log
	if roomID == "" {
		log.Println("Missing room ID") // Debug log
		http.Error(w, "Missing room ID", http.StatusBadRequest)
		return
	}

	// 3. Parse request body
	var payload struct {
		Content string `json:"content"`
	}

	bodyBytes, _ := io.ReadAll(r.Body)
	log.Printf("Raw request body: %s", string(bodyBytes)) // Debug log

	// Reset body for decoding
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("JSON decode error: %v", err) // Debug log
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Parsed content length: %d", len(payload.Content)) // Debug log

	// 4. Execute update
	result, err := db.Exec("UPDATE rooms SET content = $1 WHERE id = $2", payload.Content, roomID)
	if err != nil {
		log.Printf("Update error: %v", err) // Debug log
		http.Error(w, "Failed to save document", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Rows affected: %d", rowsAffected) // Debug log

	// 5. Success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":        "success",
		"roomId":        roomID,
		"contentLength": len(payload.Content),
	})
}

func HandleGetRoom(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	roomID := r.URL.Path[len("/api/rooms/"):]
	if roomID == "" {
		http.Error(w, "Missing room ID", http.StatusBadRequest)
		return
	}

	var content string
	err := db.QueryRow("SELECT content FROM rooms WHERE id = $1", roomID).Scan(&content)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Room not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":      roomID,
		"content": content,
	})
}
