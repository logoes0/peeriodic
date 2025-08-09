package models

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Message represents a WebSocket message
type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	UID       string    `json:"uid"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Room represents a collaborative editing room
type Room struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserUID   *string   `json:"user_uid,omitempty"` // Changed to pointer to handle NULL values
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Clients   map[*websocket.Conn]bool
	Mu        sync.RWMutex
}

// NewRoom creates a new room instance
func NewRoom(id, title string, userUID *string) *Room {
	return &Room{
		ID:        id,
		Title:     title,
		UserUID:   userUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Clients:   make(map[*websocket.Conn]bool),
	}
}

// AddClient adds a client to the room
func (r *Room) AddClient(conn *websocket.Conn) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.Clients[conn] = true
}

// RemoveClient removes a client from the room
func (r *Room) RemoveClient(conn *websocket.Conn) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	delete(r.Clients, conn)
}

// GetClientCount returns the number of clients in the room
func (r *Room) GetClientCount() int {
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	return len(r.Clients)
}

// Broadcast sends a message to all clients in the room except the sender
func (r *Room) Broadcast(message Message, sender *websocket.Conn) {
	r.Mu.RLock()
	defer r.Mu.RUnlock()

	for client := range r.Clients {
		if client != sender {
			if err := client.WriteJSON(message); err != nil {
				// Handle write error - could log or remove client
				client.Close()
				delete(r.Clients, client)
			}
		}
	}
}

// UpdateContent updates the room's content and timestamp
func (r *Room) UpdateContent(content string) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.Content = content
	r.UpdatedAt = time.Now()
}
