package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const file string = "data/data.db"

var db *sql.DB

const createDatabaseQuery string = `
  CREATE TABLE IF NOT EXISTS links (
  id INTEGER NOT NULL PRIMARY KEY,
  slug VARCHAR(255) NOT NULL,
  destination_url TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
  );`

func GetDB() (*sql.DB, error) {
	if db == nil {
		db, err := connect()
		if err != nil {
			log.Println("couldnt connect to database: " + err.Error())
		}
		return db, nil
	}
	return db, nil
}

func connect() (*sql.DB, error) {
	var err error
	db, err = sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(createDatabaseQuery); err != nil {
		return nil, err
	}
	return db, nil
}
