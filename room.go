package main

import (
	"chat/models"
	"fmt"

	"github.com/google/uuid"
)

type Room struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	Private    bool `json:"private"`
}

const welcomeMessage = "Welcome %s"

// Create new room chat
func NewRoom(name string, private bool) *Room {
	return &Room{
		ID:         uuid.New(),
		Name:       name,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
		Private:    private,
	}
}

func (room *Room) RunRoom() {
	for {
		select {

		// Register client to room
		case client := <-room.register:
			room.registerClientInRoom(client)

		// Unregister client from room
		case client := <-room.unregister:
			room.unregisterClientInRoom(client)

		// Send boardcast message to all member in room which in json format
		case message := <-room.broadcast:
			room.broadcastToClientsInRoom(message.encode())
		}
	}
}

// Member interaction
// Register and Unregister client in room
// Add client to new room and notify to all member
func (room *Room) registerClientInRoom(client *Client) {
	// If client choose public chat room
	if !room.Private {
		room.notifyClientJoined(client)
	}

	room.clients[client] = true
}

// Remove client from room
func (room *Room) unregisterClientInRoom(client *Client) {
	// If room exists remove client from that room
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
	}
}

// List online members
func (server *WsServer) listOnlineClients(client *Client) {
	for _, user := range server.users {
		message := &Message{
			Action: UserJoinedAction,
			Sender: user,
		}
		client.send <- message.encode()
	}
}

// Check if client is not yet in the room
func (client *Client) isInRoom(room *Room) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}
	return false
}

// Notification
// Notify new client joined the room
func (room *Room) notifyClientJoined(client *Client) {
	// Create welcome new member message
	message := &Message{
		Action:  SendMessageAction,
		Target:  room,
		Message: fmt.Sprintf(welcomeMessage, client.GetName()),
	}

	// Broadcast to all client in the room
	room.broadcastToClientsInRoom(message.encode())
}

// Notify client of the new room
func (client *Client) notifyRoomJoined(room *Room, sender models.User) {
	message := Message{
		Action: RoomJoinedAction,
		Target: room,
		Sender: sender,
	}

	client.send <- message.encode()
}

// Broadcast message
// Broadcast message to all members in room
func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}

// Room interaction
// Return room's ID
func (room *Room) GetID() string {
	return room.ID.String()
}

// Return room's name
func (room *Room) GetName() string {
	return room.Name
}

// Return private room or not
func (room *Room) GetPrivate() bool {
	return room.Private
}

// Find room's name in all rooms created
func (server *WsServer) FindRoomByName(name string) *Room {
	// Create loop to find all name = room[name]
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetName() == name {
			foundRoom = room
			break
		}
	}

	// If there is no room, try to create it from repo
	if foundRoom == nil {
		foundRoom = server.runRoomFormRepository(name)
	}

	return foundRoom
}

// Run room from repository
func (server *WsServer) runRoomFormRepository(name string) *Room {
	var room *Room

	// Find room
	dbRoom := server.roomRepository.FindRoomByName(name)

	// If not exists create new room
	if dbRoom != nil {
		room = NewRoom(dbRoom.GetName(), dbRoom.GetPrivate())
		room.ID, _ = uuid.Parse(dbRoom.GetID())

		go room.RunRoom()
		server.rooms[room] = true
	}

	return room
}

// Find room's ID in all rooms created
func (server *WsServer) findRoomByID(ID string) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetID() == ID {
			foundRoom = room
			break
		}
	}
	return foundRoom
}

// Create new room chat
func (server *WsServer) createRoom(name string, private bool) *Room {

	// Handle no name Room
	if name == "" {
		fmt.Println("Not such a name")
		return nil
	}

	// Make new room
	room := NewRoom(name, private)

	// Add room to repo
	server.roomRepository.AddRoom(room)

	go room.RunRoom()
	server.rooms[room] = true

	return room
}

// Join room & Leave room

// Public room
// Handle join public room action
func (client *Client) handleJoinRoomMessage(message Message) {
	// Get room name
	roomName := message.Message
	client.joinRoom(roomName, nil)
}

// Handle leave room action
func (client *Client) handleLeaveRoomMessage(message Message) {
	// Get room name
	room := client.wsServer.findRoomByID(message.Message)
	if room == nil {
		fmt.Println("Room not found")
	}

	// Delete clients
	if _, ok := client.rooms[room]; ok {
		delete(client.rooms, room)
	}

	room.unregister <- client
}

// Private room
// 	Handle join private room message
func (client *Client) handleJoinRoomPrivateMessage(message Message) {
	// Find requested client
	target := client.wsServer.findClientByID(message.Message)

	// Handle client not found
	if target == nil {
		fmt.Println("Client not found")
		return
	}

	roomName := message.Message + client.ID.String()

	// Create room for 2 clients
	client.joinRoom(roomName, target)
	target.joinRoom(roomName, client)
}

// Handle join private room
func (client *Client) joinRoom(roomName string, sender models.User) {
	room := client.wsServer.FindRoomByName(roomName)

	// Handle room not found -> Create new room
	if room == nil {
		fmt.Println("Create new room : " + roomName)
		room = client.wsServer.createRoom(roomName, sender != nil)
	}

	// Don't allow to join private rooms through public room
	if sender == nil && room.Private {
		fmt.Println("Not allow to join private room through public room")
		return
	}

	// If there is client in room, register client to room and notify all member
	if !client.isInRoom(room) {
		client.rooms[room] = true
		room.register <- client
		client.notifyRoomJoined(room, sender)
	}
}
