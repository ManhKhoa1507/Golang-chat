package config

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (

	// Define create table action
	CreateRoomTable = `
	CREATE TABLE IF NOT EXISTS room (
		id VARCHAR (25) NOT NULL PRIMARY KEY,
		name VARCHAR(25) NOT NULL,
		private TINYINT NULL
	);
	`

	CreateUserTable = `
	CREATE TABLE IF NOT EXISTS user(
		id VARCHAR(25) NOT NULL PRIMARY KEY,
		name VARCHAR(25) NOT NULL,
	);
	`
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
	_, err = db.Exec(CreateRoomTable)

	// Handle error
	if err != nil {
		fmt.Println("Error when creating table room")
	}

	_, err = db.Exec(CreateUserTable)
	return db
}
