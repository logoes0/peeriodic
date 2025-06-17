package routes

import (
	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/logoes0/peeriodic.git/server"
)

func SetupRoutes(app *fiber.App, authClient *auth.Client) {

	// middleware to upgrade only WebSocket requests
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			// Store token in locals (not entire context)
			c.Locals("token", c.Query("token"))
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(server.HandleWebSocket(authClient)))
}
