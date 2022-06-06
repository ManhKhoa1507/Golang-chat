package repository

import (
	"chat/models"
	"database/sql"
	"fmt"
)

const (
	InsertRoom = "INSERT INTO room(id, name, private) values (?, ?, ?)"
	FindRoom   = "SELECT id, name, private FROM room WHERE name = ? LIMIT 1"
)

type Room struct {
	ID      string
	Name    string
	Private bool
}

// Get room's ID
func (room *Room) GetID() string {
	return room.ID
}

// Get room's name
func (room *Room) GetName() string {
	return room.Name
}

// Get room's Private
func (room *Room) GetPrivate() bool {
	return room.Private
}

type RoomRepository struct {
	Db *sql.DB
}

func HandleError(err error, nameError string) {
	if err != nil {
		fmt.Println("Error ", nameError)
	}
}

// Add value to room table
func (repo *RoomRepository) AddRoom(room models.Room) {
	stmt, err := repo.Db.Prepare(InsertRoom)

	// Add roomID, roomName, roomPrivate
	_, err = stmt.Exec(room.GetID(), room.GetName(), room.GetPrivate())
	// Handle error
	HandleError(err, "Insert value to Room")
}

func (repo *RoomRepository) findRoomByName(name string) models.Room {
	row := repo.Db.QueryRow(FindRoom, name)

	var room Room

	// Handle error
	if err := row.Scan(&room.ID, &room.Name, &room.Private); err != nil {
		// Room not found
		if err == sql.ErrNoRows {
			fmt.Println("Room not found")
			return nil
		}
	}

	return &room
}
