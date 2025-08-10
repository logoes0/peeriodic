package routers

import (
	"log"
	"net/http"
	"strings"

	"github.com/logoes0/peeriodic.git/handlers"
	"github.com/logoes0/peeriodic.git/middleware"
	"github.com/logoes0/peeriodic.git/services"
)

// Router handles all HTTP routing
type Router struct {
	roomHandler     *handlers.RoomHandler
	documentHandler *handlers.DocumentHandler
	wsService       *services.WebSocketService
}

// NewRouter creates a new router instance
func NewRouter(dbService *services.DatabaseService, wsService *services.WebSocketService) *Router {
	return &Router{
		roomHandler:     handlers.NewRoomHandler(dbService),
		documentHandler: handlers.NewDocumentHandler(dbService),
		wsService:       wsService,
	}
}

// SetupRoutes configures all application routes with middleware
func (r *Router) SetupRoutes() {
	// WebSocket endpoint - NO middleware (WebSocket needs direct access to response writer)
	http.HandleFunc("/ws", r.handleWebSocket)

	// HTTP API endpoints - apply CORS middleware
	http.HandleFunc("/api/rooms", middleware.Logging(middleware.CORS(r.handleRooms)))
	http.HandleFunc("/api/save", middleware.Logging(middleware.CORS(r.handleSave)))

	// Handle room-specific operations with path parameters
	http.HandleFunc("/api/rooms/", middleware.Logging(middleware.CORS(r.handleRoomOperations)))
}

// handleWebSocket handles WebSocket connections
func (r *Router) handleWebSocket(w http.ResponseWriter, req *http.Request) {
	r.wsService.HandleConnection(w, req, r.roomHandler.DBService)
}

// handleRooms handles room listing and creation
func (r *Router) handleRooms(w http.ResponseWriter, req *http.Request) {
	r.roomHandler.HandleRooms(w, req)
}

// handleRoomOperations handles room-specific operations (GET, DELETE)
func (r *Router) handleRoomOperations(w http.ResponseWriter, req *http.Request) {
	log.Printf("handleRoomOperations called with path: %s", req.URL.Path)

	// Extract room ID from path
	pathParts := strings.Split(req.URL.Path, "/")
	log.Printf("Path parts: %v", pathParts)

	if len(pathParts) < 4 {
		log.Printf("Invalid path length: %d", len(pathParts))
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}
	roomID := pathParts[3]
	log.Printf("Extracted room ID: %s", roomID)

	if roomID == "" {
		log.Printf("Empty room ID")
		http.Error(w, "Missing room ID", http.StatusBadRequest)
		return
	}

	// Create a new request with the room ID in the path for the handlers
	req.URL.Path = "/api/rooms/" + roomID
	log.Printf("Modified path: %s", req.URL.Path)

	switch req.Method {
	case http.MethodGet:
		log.Printf("Handling GET request for room: %s", roomID)
		r.roomHandler.HandleRoomByID(w, req)
	case http.MethodDelete:
		log.Printf("Handling DELETE request for room: %s", roomID)
		r.roomHandler.HandleDeleteRoom(w, req)
	default:
		log.Printf("Method not allowed: %s", req.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleSave handles document saving
func (r *Router) handleSave(w http.ResponseWriter, req *http.Request) {
	r.documentHandler.HandleSave(w, req)
}
