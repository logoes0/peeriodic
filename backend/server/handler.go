package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func HandleViewerWebSocket() func(*websocket.Conn) {
	return func(conn *websocket.Conn) {
		ctxRaw := conn.Locals("ctx")
		if ctxRaw == nil {
			log.Println("Missing context for viewer")
			conn.Close()
			return
		}
		ctx := ctxRaw.(*fiber.Ctx)
		sessionID := ctx.Params("sessionId")
		log.Printf("📺 Viewer connected to session: %s", sessionID)

		vMu.Lock()
		viewers[sessionID] = append(viewers[sessionID], conn)
		vMu.Unlock()

		defer func() {
			vMu.Lock()
			conns := viewers[sessionID]
			for i, c := range conns {
				if c == conn {
					viewers[sessionID] = append(conns[:i], conns[i+1:]...)
					break
				}
			}
			vMu.Unlock()
			conn.Close()
		}()

		// viewer doesn't send anything, just holds the connection open
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}
}
