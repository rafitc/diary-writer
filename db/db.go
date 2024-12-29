package db

import (
	"database/sql"
	"main/logger"

	_ "github.com/mattn/go-sqlite3"
)

var log = logger.Logger

type Database struct {
	Conn *sql.DB
}

func NewDatabase(dataSourceName string) (*Database, error) {
	log.Debugf("Trying to connect to db %s", dataSourceName)
	conn, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return &Database{Conn: conn}, nil
}

func (db *Database) Fetch(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (db *Database) Insert(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.Conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *Database) Update(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.Conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *Database) Delete(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.Conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *Database) Close() error {
	return db.Conn.Close()
}
