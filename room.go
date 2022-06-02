package main

import "fmt"

type Room struct {
	// Define Room type struct
	name       string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

const welcomeMessage = "Welcome %s"

// Create new room chat
func NewRoom(name string) *Room {
	return &Room{
		name:       name,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
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

// Notify new client joined the room
func (room *Room) notifyClientJoined(client *Client) {
	// Create welcome new member message
	message := &Message{
		Action:  SendMessageAction,
		Target:  room.name,
		Message: fmt.Sprintf(welcomeMessage, client.getName()),
	}

	// Broadcast to all client in the room
	room.broadcastToClientsInRoom(message.encode())
}

// Add client to new room and notify to all member
func (room *Room) registerClientInRoom(client *Client) {
	room.notifyClientJoined(client)
	room.clients[client] = true
}

// Remove client from room
func (room *Room) unregisterClientInRoom(client *Client) {
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
	}
}

// Broadcast message to all members in room
func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}

func (room *Room) GetName() string {
	return room.name
}

// Find room's name in all rooms created
func (server *WsServer) findRoomByName(name string) *Room {
	// Create loop to find all name = room[name]
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetName() == name {
			foundRoom = room
			break
		}
	}
	return foundRoom
}

// Create new room chat
func (server *WsServer) createRoom(name string) *Room {
	room := NewRoom(name)
	go room.RunRoom()
	server.rooms[room] = true
	return room
}

// Handle join room action
func (client *Client) handleJoinRoomMessage(message Message) {
	// Get room name
	roomName := message.Message

	// Find room by name
	room := client.wsServer.findRoomByName(roomName)

	// If room not exists -> Create new room
	if room == nil {
		room = client.wsServer.createRoom(roomName)
	}

	// Add client to new room
	client.rooms[room] = true
	room.register <- client
}

// Handle leave room action
func (client *Client) handleLeaveRoomMessage(message Message) {
	// Get room name
	roomName := message.Message

	// Find room and delete client from that room
	room := client.wsServer.findRoomByName(roomName)
	if _, ok := client.rooms[room]; ok {
		delete(client.rooms, room)
	}

	// unregister client from room
	room.unregister <- client
}
