-- Open SQLite3 command line or a script

-- Create a new database (if you haven't already)
sqlite3 diary.db

-- Create a table named 'daily_updates' with the specified structure
CREATE TABLE daily_updates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    asset TEXT,
    asset_extension TEXT,
    asset_download_link TEXT,
    asset_blob BLOB,
    creation_date DATE DEFAULT CURRENT_DATE,
    is_updated BOOLEAN DEFAULT 0
);

-- Exit the SQLite3 command line
.exit
