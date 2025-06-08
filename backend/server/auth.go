package server

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
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
		defer conn.Close()

		// Defensive check
		tokenRaw := conn.Locals("token")
		token, ok := tokenRaw.(string)
		if !ok || token == "" {
			fmt.Println("❌ Missing token")
			return
		}

		decoded, err := authClient.VerifyIDToken(context.Background(), token)
		if err != nil {
			fmt.Printf("❌ Invalid token: %v\n", err)
			return
		}

		uid := decoded.UID
		email := decoded.Claims["email"]
		fmt.Printf("✅ Authenticated UID: %s, Email: %v\n", uid, email)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("WebSocket read error: %v\n", err)
				break
			}

			fmt.Printf("From %s: %s\n", uid, msg)

			// Echo "received"
			if err := conn.WriteMessage(websocket.TextMessage, []byte("received")); err != nil {
				fmt.Printf("WebSocket write error: %v\n", err)
				break
			}
		}
	}
}
