-- Open SQLite3 command line or a script

-- Create a new database (if you haven't already)
sqlite3 diary.db

-- Create a table named 'my_table' with the specified structure
CREATE TABLE daily_updates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    asset TEXT,
    creation_date DATE DEFAULT CURRENT_DATE,
    is_updated BOOLEAN DEFAULT 0
);

-- Exit the SQLite3 command line
.exit