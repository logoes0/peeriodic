package services

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/logoes0/peeriodic.git/config"
	"github.com/logoes0/peeriodic.git/models"
)

// WebSocketService handles WebSocket connections and real-time communication
type WebSocketService struct {
	config *config.Config
	rooms  map[string]*RoomManager
	mu     sync.RWMutex
}

// RoomManager manages clients and document state for a specific room
type RoomManager struct {
	ID       string
	Clients  map[*websocket.Conn]bool
	Document string
	mu       sync.RWMutex
}

// NewWebSocketService creates a new WebSocket service instance
func NewWebSocketService(cfg *config.Config) *WebSocketService {
	return &WebSocketService{
		config: cfg,
		rooms:  make(map[string]*RoomManager),
	}
}

// GetUpgrader returns a configured WebSocket upgrader
func (ws *WebSocketService) GetUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  ws.config.WebSocket.ReadBufferSize,
		WriteBufferSize: ws.config.WebSocket.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return ws.config.WebSocket.CheckOrigin
		},
	}
}

// HandleConnection handles a new WebSocket connection
func (ws *WebSocketService) HandleConnection(w http.ResponseWriter, r *http.Request, dbService *DatabaseService) {
	// Extract room ID from query parameters
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		http.Error(w, "Missing room ID", http.StatusBadRequest)
		return
	}

	// Ensure room exists in database
	room, err := dbService.EnsureRoomExists(roomID)
	if err != nil {
		log.Printf("Failed to ensure room exists: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Upgrade HTTP connection to WebSocket
	upgrader := ws.GetUpgrader()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer ws.closeConnection(conn, roomID)

	// Get or create room manager
	roomManager := ws.getOrCreateRoom(roomID, room.Content)

	// Add client to room
	roomManager.mu.Lock()
	roomManager.Clients[conn] = true
	clientCount := len(roomManager.Clients)
	roomManager.mu.Unlock()

	log.Printf("Client connected to room %s (total clients: %d)", roomID, clientCount)

	// Send initial document state to new client
	if err := conn.WriteJSON(models.Message{
		Type: "init",
		Data: roomManager.Document,
	}); err != nil {
		log.Printf("Failed to send initial document: %v", err)
		return
	}

	// Handle incoming messages
	ws.handleMessages(conn, roomManager, dbService)
}

// handleMessages processes incoming WebSocket messages
func (ws *WebSocketService) handleMessages(conn *websocket.Conn, roomManager *RoomManager, dbService *DatabaseService) {
	for {
		var msg models.Message
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Client disconnected unexpectedly: %v", err)
			}
			break
		}

		switch msg.Type {
		case "update":
			ws.handleDocumentUpdate(conn, roomManager, msg.Data, dbService)
		default:
			log.Printf("Unknown message type: %s", msg.Type)
		}
	}
}

// handleDocumentUpdate processes document update messages
func (ws *WebSocketService) handleDocumentUpdate(conn *websocket.Conn, roomManager *RoomManager, content string, dbService *DatabaseService) {
	roomManager.mu.Lock()
	defer roomManager.mu.Unlock()

	// Update local document state
	roomManager.Document = content

	// Broadcast to other clients in the room
	for client := range roomManager.Clients {
		if client != conn {
			if err := client.WriteJSON(models.Message{
				Type: "update",
				Data: content,
			}); err != nil {
				log.Printf("Failed to broadcast to client: %v", err)
				client.Close()
				delete(roomManager.Clients, client)
			}
		}
	}

	// Persist to database asynchronously
	go func() {
		log.Printf("Persisting document update for room %s, content length: %d", roomManager.ID, len(content))
		if err := dbService.UpdateRoomContent(roomManager.ID, content); err != nil {
			log.Printf("Failed to persist document update: %v", err)
		} else {
			log.Printf("Successfully persisted document update for room %s", roomManager.ID)
		}
	}()
}

// getOrCreateRoom returns an existing room manager or creates a new one
func (ws *WebSocketService) getOrCreateRoom(roomID, initialContent string) *RoomManager {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if room, exists := ws.rooms[roomID]; exists {
		return room
	}

	room := &RoomManager{
		ID:       roomID,
		Clients:  make(map[*websocket.Conn]bool),
		Document: initialContent,
	}
	ws.rooms[roomID] = room
	return room
}

// closeConnection removes a client from the room and cleans up if necessary
func (ws *WebSocketService) closeConnection(conn *websocket.Conn, roomID string) {
	conn.Close()

	ws.mu.RLock()
	room, exists := ws.rooms[roomID]
	ws.mu.RUnlock()

	if !exists {
		return
	}

	room.mu.Lock()
	delete(room.Clients, conn)
	remainingClients := len(room.Clients)
	room.mu.Unlock()

	log.Printf("Client disconnected from room %s (%d clients remaining)", roomID, remainingClients)

	// Clean up empty rooms
	if remainingClients == 0 {
		ws.mu.Lock()
		delete(ws.rooms, roomID)
		ws.mu.Unlock()
		log.Printf("Room %s cleaned up (no clients remaining)", roomID)
	}
}

// BroadcastToRoom sends a message to all clients in a specific room
func (ws *WebSocketService) BroadcastToRoom(roomID string, message models.Message) {
	ws.mu.RLock()
	room, exists := ws.rooms[roomID]
	ws.mu.RUnlock()

	if !exists {
		return
	}

	room.mu.RLock()
	defer room.mu.RUnlock()

	for client := range room.Clients {
		if err := client.WriteJSON(message); err != nil {
			log.Printf("Failed to broadcast to client: %v", err)
			client.Close()
			delete(room.Clients, client)
		}
	}
}

// GetRoomStats returns statistics about active rooms
func (ws *WebSocketService) GetRoomStats() map[string]int {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	stats := make(map[string]int)
	for roomID, room := range ws.rooms {
		room.mu.RLock()
		stats[roomID] = len(room.Clients)
		room.mu.RUnlock()
	}
	return stats
}
