package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/logoes0/peeriodic.git/config"
)

// DatabaseService handles all database operations
type DatabaseService struct {
	db *sql.DB
}

// NewDatabaseService creates a new database service instance
func NewDatabaseService(cfg *config.Config) (*DatabaseService, error) {
	db, err := sql.Open("postgres", cfg.GetDatabaseConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connection established successfully")
	return &DatabaseService{db: db}, nil
}

// Close closes the database connection
func (ds *DatabaseService) Close() error {
	return ds.db.Close()
}

// User represents a user in the database
type User struct {
	ID        int       `json:"id"`
	UID       string    `json:"uid"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Room represents a room in the database
type Room struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserUID   *string   `json:"user_uid"` // Changed to pointer to handle NULL values
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUser creates a new user in the database
func (ds *DatabaseService) CreateUser(uid, email, name string) (*User, error) {
	query := `
		INSERT INTO users (uid, email, name, created_at, updated_at) 
		VALUES ($1, $2, $3, NOW(), NOW()) 
		RETURNING id, uid, email, name, created_at, updated_at
	`

	user := &User{}
	err := ds.db.QueryRow(query, uid, email, name).Scan(
		&user.ID, &user.UID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByUID retrieves a user by UID
func (ds *DatabaseService) GetUserByUID(uid string) (*User, error) {
	query := `
		SELECT id, uid, email, name, created_at, updated_at 
		FROM users 
		WHERE uid = $1
	`

	user := &User{}
	err := ds.db.QueryRow(query, uid).Scan(
		&user.ID, &user.UID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %s", uid)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// EnsureUserExists creates a user if it doesn't exist, otherwise returns the existing user
func (ds *DatabaseService) EnsureUserExists(uid, email, name string) (*User, error) {
	query := `
		INSERT INTO users (uid, email, name, created_at, updated_at) 
		VALUES ($1, $2, $3, NOW(), NOW()) 
		ON CONFLICT (uid) DO UPDATE SET 
			email = EXCLUDED.email,
			name = EXCLUDED.name,
			updated_at = NOW()
		RETURNING id, uid, email, name, created_at, updated_at
	`

	user := &User{}
	err := ds.db.QueryRow(query, uid, email, name).Scan(
		&user.ID, &user.UID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure user exists: %w", err)
	}

	return user, nil
}

// CreateRoom creates a new room in the database
func (ds *DatabaseService) CreateRoom(id, title string, userUID *string) (*Room, error) {
	query := `
		INSERT INTO rooms (id, title, user_uid, content, created_at, updated_at) 
		VALUES ($1, $2, $3, '', NOW(), NOW()) 
		RETURNING id, title, content, user_uid, created_at, updated_at
	`

	room := &Room{}
	err := ds.db.QueryRow(query, id, title, userUID).Scan(
		&room.ID, &room.Title, &room.Content, &room.UserUID, &room.CreatedAt, &room.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}

	return room, nil
}

// GetRoom retrieves a room by ID
func (ds *DatabaseService) GetRoom(id string) (*Room, error) {
	query := `
		SELECT id, title, content, user_uid, created_at, updated_at 
		FROM rooms 
		WHERE id = $1
	`

	room := &Room{}
	err := ds.db.QueryRow(query, id).Scan(
		&room.ID, &room.Title, &room.Content, &room.UserUID, &room.CreatedAt, &room.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("room not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	return room, nil
}

// GetRoomsByUser retrieves all rooms for a specific user
func (ds *DatabaseService) GetRoomsByUser(userUID string) ([]*Room, error) {
	query := `
		SELECT id, title, content, user_uid, created_at, updated_at 
		FROM rooms 
		WHERE user_uid = $1 
		ORDER BY updated_at DESC
	`

	rows, err := ds.db.Query(query, userUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rooms: %w", err)
	}
	defer rows.Close()

	var rooms []*Room
	for rows.Next() {
		room := &Room{}
		err := rows.Scan(
			&room.ID, &room.Title, &room.Content, &room.UserUID, &room.CreatedAt, &room.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan room: %w", err)
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

// UpdateRoomContent updates the content of a room
func (ds *DatabaseService) UpdateRoomContent(id, content string) error {
	query := `
		UPDATE rooms 
		SET content = $1, updated_at = NOW() 
		WHERE id = $2
	`

	log.Printf("Executing update query for room %s with content length %d", id, len(content))
	result, err := ds.db.Exec(query, content, id)
	if err != nil {
		log.Printf("Database error updating room %s: %v", id, err)
		return fmt.Errorf("failed to update room content: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for room %s: %v", id, err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	log.Printf("Updated %d rows for room %s", rowsAffected, id)
	if rowsAffected == 0 {
		log.Printf("No rows affected for room %s - room may not exist", id)
		return fmt.Errorf("room not found: %s", id)
	}

	return nil
}

// DeleteRoom deletes a room by ID
func (ds *DatabaseService) DeleteRoom(id string) error {
	query := `DELETE FROM rooms WHERE id = $1`

	result, err := ds.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("room not found: %s", id)
	}

	return nil
}

// EnsureRoomExists creates a room if it doesn't exist, otherwise returns the existing room
func (ds *DatabaseService) EnsureRoomExists(id string) (*Room, error) {
	query := `
		INSERT INTO rooms (id, title, content, user_uid, created_at, updated_at) 
		VALUES ($1, 'Untitled Room', '', NULL, NOW(), NOW()) 
		ON CONFLICT (id) DO UPDATE SET id = EXCLUDED.id 
		RETURNING id, title, content, user_uid, created_at, updated_at
	`

	log.Printf("Ensuring room exists: %s", id)
	room := &Room{}
	err := ds.db.QueryRow(query, id).Scan(
		&room.ID, &room.Title, &room.Content, &room.UserUID, &room.CreatedAt, &room.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to ensure room exists for %s: %v", id, err)
		return nil, fmt.Errorf("failed to ensure room exists: %w", err)
	}

	log.Printf("Room %s exists with content length: %d", id, len(room.Content))
	return room, nil
}
