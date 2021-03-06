package models

// Add rooms interfaces

type Room interface {
	GetID() string
	GetName() string
	GetPrivate() bool
}

type RoomRepository interface {
	AddRoom(room Room)
	FindRoomByName(name string) Room
}
