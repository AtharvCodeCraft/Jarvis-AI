import sqlite3

try:
    conn = sqlite3.connect('jarvis.db')
    cursor = conn.cursor()
    cursor.execute("SELECT name FROM sqlite_master WHERE type='table'")
    tables = cursor.fetchall()
    print("Tables in DB:", tables)
    
    if ('users',) in tables:
        cursor.execute("PRAGMA table_info(users)")
        print("Users schema:", cursor.fetchall())
    else:
        print("TABLE 'users' DOES NOT EXIST!")
except Exception as e:
    print("DB Error:", e)
