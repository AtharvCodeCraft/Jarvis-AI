package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite", "./jarvis.db")
	if err != nil {
		return err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS memory (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT UNIQUE NOT NULL,
		value TEXT NOT NULL
	);`

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		return err
	}

	log.Println("SQLite database initialized successfully.")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
