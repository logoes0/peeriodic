package server

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

var (
	viewers = make(map[string][]*websocket.Conn)
	vMu     sync.Mutex
)
