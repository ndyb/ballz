package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Highscore represents a row in the highscores table
type Highscore struct {
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}

// DB is a wrapper around sql.DB
type DB struct {
	*sql.DB
}

// Connect establishes a connection to the database
func Connect() (*DB, error) {
	// Get connection details from environment variables
	host := os.Getenv("SUPABASE_HOST")
	port := os.Getenv("SUPABASE_PORT")
	user := os.Getenv("SUPABASE_USER")
	password := os.Getenv("SUPABASE_PASSWORD")
	dbname := os.Getenv("SUPABASE_DBNAME")

	// Construct connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to database")
	return &DB{db}, nil
}

// GetHighscores retrieves highscores ordered by score desc with a limit
func (db *DB) GetHighscores(limit int) ([]Highscore, error) {
	rows, err := db.Query("SELECT score, created_at FROM highscores ORDER BY score DESC LIMIT $1", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var highscores []Highscore
	for rows.Next() {
		var h Highscore
		if err := rows.Scan(&h.Score, &h.CreatedAt); err != nil {
			return nil, err
		}
		highscores = append(highscores, h)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return highscores, nil
}

// InsertHighscore adds a new highscore to the table
// Note: created_at is set automatically by PostgreSQL
func (db *DB) InsertHighscore(score int) error {
	_, err := db.Exec("INSERT INTO highscores (score) VALUES ($1)", score)
	return err
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
