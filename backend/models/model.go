package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type Room struct {
	Clients  map[*websocket.Conn]bool
	Document string
	Mu       sync.Mutex
}
