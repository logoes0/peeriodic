package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/logoes0/peeriodic.git/services"
	"github.com/logoes0/peeriodic.git/utils"
)

// RoomHandler handles room-related HTTP requests
type RoomHandler struct {
	DBService *services.DatabaseService
}

// NewRoomHandler creates a new room handler instance
func NewRoomHandler(dbService *services.DatabaseService) *RoomHandler {
	return &RoomHandler{
		DBService: dbService,
	}
}

// CreateRoomRequest represents the request body for creating a room
type CreateRoomRequest struct {
	Title string `json:"title"`
	UID   string `json:"uid"`
}

// RoomResponse represents a room in API responses
type RoomResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// HandleRooms handles room listing and creation
func (rh *RoomHandler) HandleRooms(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rh.handleGetRooms(w, r)
	case http.MethodPost:
		rh.handleCreateRoom(w, r)
	default:
		utils.MethodNotAllowed(w)
	}
}

// handleGetRooms retrieves rooms for a specific user
func (rh *RoomHandler) handleGetRooms(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	if uid == "" {
		utils.BadRequest(w, "Missing uid parameter")
		return
	}

	rooms, err := rh.DBService.GetRoomsByUser(uid)
	if err != nil {
		utils.InternalServerError(w, "Failed to retrieve rooms")
		return
	}

	// Convert to response format
	var response []RoomResponse
	for _, room := range rooms {
		response = append(response, RoomResponse{
			ID:        room.ID,
			Title:     room.Title,
			CreatedAt: room.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	utils.SuccessResponse(w, response)
}

// handleCreateRoom creates a new room
func (rh *RoomHandler) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	var req CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	if req.UID == "" {
		utils.BadRequest(w, "UID is required")
		return
	}

	if req.Title == "" {
		req.Title = "Untitled Room"
	}

	roomID := uuid.New().String()
	room, err := rh.DBService.CreateRoom(roomID, req.Title, req.UID)
	if err != nil {
		utils.InternalServerError(w, "Failed to create room")
		return
	}

	response := RoomResponse{
		ID:    room.ID,
		Title: room.Title,
	}

	utils.SuccessResponse(w, response)
}

// HandleRoomByID handles getting a specific room by ID
func (rh *RoomHandler) HandleRoomByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.MethodNotAllowed(w)
		return
	}

	// Extract room ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		utils.BadRequest(w, "Invalid room ID")
		return
	}
	roomID := pathParts[3]

	if roomID == "" {
		utils.BadRequest(w, "Missing room ID")
		return
	}

	room, err := rh.DBService.GetRoom(roomID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.NotFound(w, "Room not found")
		} else {
			utils.InternalServerError(w, "Failed to retrieve room")
		}
		return
	}

	response := RoomResponse{
		ID:      room.ID,
		Title:   room.Title,
		Content: room.Content,
	}

	utils.SuccessResponse(w, response)
}

// HandleDeleteRoom handles room deletion
func (rh *RoomHandler) HandleDeleteRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.MethodNotAllowed(w)
		return
	}

	// Extract room ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		utils.BadRequest(w, "Invalid room ID")
		return
	}
	roomID := pathParts[3]

	if roomID == "" {
		utils.BadRequest(w, "Missing room ID")
		return
	}

	err := rh.DBService.DeleteRoom(roomID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.NotFound(w, "Room not found")
		} else {
			utils.InternalServerError(w, "Failed to delete room")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
