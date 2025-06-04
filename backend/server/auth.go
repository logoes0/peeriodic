package server

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"

	"google.golang.org/api/option"
)

func InitFirebaseAuth() (*auth.Client, error) {
	_ = godotenv.Load()

	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing Firebase app: %v", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %v", err)
	}

	return authClient, nil
}

func HandleWebSocket(authClient *auth.Client) func(*websocket.Conn) {
	return func(conn *websocket.Conn) {
		// Access original HTTP context
		ctxRaw := conn.Locals("ctx")
		if ctxRaw == nil {
			log.Println("Missing HTTP context")
			conn.Close()
			return
		}

		ctx := ctxRaw.(*fiber.Ctx)

		// Extract token from query params
		token := ctx.Query("token")
		if token == "" {
			log.Println("Missing token in query string")
			conn.Close()
			return
		}

		// Verify Firebase ID token
		decoded, err := authClient.VerifyIDToken(context.Background(), token)
		if err != nil {
			log.Printf("Invalid Firebase token: %v", err)
			conn.Close()
			return
		}

		// Extract user info
		uid := decoded.UID
		email := decoded.Claims["email"]
		log.Printf("✅ Authenticated UID: %s, Email: %v", uid, email)

		// Now start handling WebSocket messages
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				break
			}
			log.Printf("From %s: %s", uid, msg)

			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("WebSocket write error: %v", err)
				break
			}
		}
	}
}
