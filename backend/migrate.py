import sqlite3
import os

db_path = "jarvis.db"

if os.path.exists(db_path):
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    
    # Check if user_id column exists in conversations
    cursor.execute("PRAGMA table_info(conversations)")
    columns = [col[1] for col in cursor.fetchall()]
    
    if "user_id" not in columns:
        print("Adding user_id to conversations table")
        cursor.execute("ALTER TABLE conversations ADD COLUMN user_id VARCHAR NOT NULL DEFAULT 'default_user'")
        
    conn.commit()
    conn.close()
    print("Migration complete.")
else:
    print("Database not found, no migration needed.")
