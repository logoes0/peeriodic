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

// Room represents a room in the database
type Room struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserUID   string    `json:"user_uid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateRoom creates a new room in the database
func (ds *DatabaseService) CreateRoom(id, title, userUID string) (*Room, error) {
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

	result, err := ds.db.Exec(query, content, id)
	if err != nil {
		return fmt.Errorf("failed to update room content: %w", err)
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
		INSERT INTO rooms (id, title, content, created_at, updated_at) 
		VALUES ($1, 'Untitled Room', '', NOW(), NOW()) 
		ON CONFLICT (id) DO UPDATE SET id = EXCLUDED.id 
		RETURNING id, title, content, user_uid, created_at, updated_at
	`

	room := &Room{}
	err := ds.db.QueryRow(query, id).Scan(
		&room.ID, &room.Title, &room.Content, &room.UserUID, &room.CreatedAt, &room.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure room exists: %w", err)
	}

	return room, nil
}

