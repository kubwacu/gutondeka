package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Initialize() error {
	dbPath := filepath.Join("database.db")
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	initSQL, err := os.ReadFile(filepath.Join("internal", "db", "migrations", "init.sql"))
	if err != nil {
		return err
	}

	_, err = db.Exec(string(initSQL))
	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

// CheckFileExists checks if an article with the same filename exists
func CheckFileExists(filename string) (bool, error) {
	filePath := filepath.Join("data", "articles", filename)

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM articles WHERE file_path = ?)`
	err := db.QueryRow(query, filePath).Scan(&exists)

	return exists, err
}
