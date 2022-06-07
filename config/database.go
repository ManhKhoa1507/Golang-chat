package config

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	createRoom = "CREATE TABLE IF NOT EXISTS room (id VARCHAR(255) NOT NULL PRIMARY KEY, name VARCHAR(255) NOT NULL, private TINYINT NULL);"
	createUser = "CREATE TABLE IF NOT EXISTS user (id VARCHAR(25) NOT NULL PRIMARY KEY, name VARCHAR(25) NOT NULL);"
)

// Init the databse
func InitDB() *sql.DB {
	// Open database file
	db, err := sql.Open("sqlite3", "./chatdb.db")

	// Handle error
	if err != nil {
		panic("Error when open database")
	}

	// Create database by using MySQl

	// Create database
	createTableRoomQuery, err := db.Prepare(createRoom)
	// Handle error
	if err != nil {
		fmt.Println("Error when creating table room")
	}
	createTableRoomQuery.Exec()

	// Create database
	createTableUserQuery, err := db.Prepare(createUser)
	// Handle error
	if err != nil {
		fmt.Println("Error when creating table room")
	}
	createTableUserQuery.Exec()
	return db
}
