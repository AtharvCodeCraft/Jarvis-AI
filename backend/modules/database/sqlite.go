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

	// Create legacy memory table first (to prevent breaking legacy key-value calls if any)
	legacyTable := `
	CREATE TABLE IF NOT EXISTS memory (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT UNIQUE NOT NULL,
		value TEXT NOT NULL
	);`
	if _, err = DB.Exec(legacyTable); err != nil {
		return err
	}

	// Schema modifications for production Agentic platform
	schemas := []string{
		`CREATE TABLE IF NOT EXISTS conversations (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			conversation_id TEXT REFERENCES conversations(id) ON DELETE CASCADE,
			sender TEXT CHECK(sender IN ('user', 'assistant', 'system')),
			content TEXT NOT NULL,
			agent_assigned TEXT,
			tokens_used INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS memory_nodes (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			memory_type TEXT CHECK(memory_type IN ('semantic', 'preference', 'episodic')),
			content TEXT NOT NULL,
			last_accessed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS documents (
			id TEXT PRIMARY KEY,
			filename TEXT NOT NULL,
			file_type TEXT NOT NULL,
			file_path TEXT NOT NULL,
			checksum TEXT UNIQUE NOT NULL,
			chunk_count INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS tool_logs (
			id TEXT PRIMARY KEY,
			agent_name TEXT NOT NULL,
			tool_name TEXT NOT NULL,
			arguments TEXT NOT NULL,
			execution_status TEXT CHECK(execution_status IN ('pending', 'approved', 'rejected', 'success', 'failed')),
			error_message TEXT,
			duration_ms INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, schema := range schemas {
		if _, err = DB.Exec(schema); err != nil {
			return err
		}
	}

	log.Println("Advanced database schemas initialized successfully.")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
