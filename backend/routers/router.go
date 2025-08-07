package routers

import (
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
	// Apply middleware to all routes
	http.HandleFunc("/ws", middleware.Logging(middleware.CORS(r.handleWebSocket)))
	http.HandleFunc("/api/rooms", middleware.Logging(middleware.CORS(r.handleRooms)))
	http.HandleFunc("/api/rooms/", middleware.Logging(middleware.CORS(r.handleRoomOperations)))
	http.HandleFunc("/api/save", middleware.Logging(middleware.CORS(r.handleSave)))
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
	// Extract room ID from path
	pathParts := strings.Split(req.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}
	roomID := pathParts[3]

	if roomID == "" {
		http.Error(w, "Missing room ID", http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodGet:
		r.roomHandler.HandleRoomByID(w, req)
	case http.MethodDelete:
		r.roomHandler.HandleDeleteRoom(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleSave handles document saving
func (r *Router) handleSave(w http.ResponseWriter, req *http.Request) {
	r.documentHandler.HandleSave(w, req)
}
