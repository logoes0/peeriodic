package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/logoes0/peeriodic.git/services"
	"github.com/logoes0/peeriodic.git/utils"
)

// DocumentHandler handles document-related HTTP requests
type DocumentHandler struct {
	dbService *services.DatabaseService
}

// NewDocumentHandler creates a new document handler instance
func NewDocumentHandler(dbService *services.DatabaseService) *DocumentHandler {
	return &DocumentHandler{
		dbService: dbService,
	}
}

// SaveDocumentRequest represents the request body for saving a document
type SaveDocumentRequest struct {
	Content string `json:"content"`
}

// SaveDocumentResponse represents the response for saving a document
type SaveDocumentResponse struct {
	Status        string `json:"status"`
	RoomID        string `json:"roomId"`
	ContentLength int    `json:"contentLength"`
}

// HandleSave handles document saving
func (dh *DocumentHandler) HandleSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.MethodNotAllowed(w)
		return
	}

	// Get room ID from query parameters
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		utils.BadRequest(w, "Missing room ID")
		return
	}

	// Parse request body
	var req SaveDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	// Validate content
	if req.Content == "" {
		utils.BadRequest(w, "Content cannot be empty")
		return
	}

	// Save document to database
	err := dh.dbService.UpdateRoomContent(roomID, req.Content)
	if err != nil {
		utils.InternalServerError(w, "Failed to save document")
		return
	}

	response := SaveDocumentResponse{
		Status:        "success",
		RoomID:        roomID,
		ContentLength: len(req.Content),
	}

	utils.SuccessResponse(w, response)
}

// HandleGetDocument handles retrieving a document
func (dh *DocumentHandler) HandleGetDocument(w http.ResponseWriter, r *http.Request) {
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

	room, err := dh.dbService.GetRoom(roomID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.NotFound(w, "Room not found")
		} else {
			utils.InternalServerError(w, "Failed to retrieve document")
		}
		return
	}

	response := map[string]string{
		"id":      room.ID,
		"content": room.Content,
	}

	utils.SuccessResponse(w, response)
}

