package main

type Room struct {
	// Define Room type struct
	name       string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

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

		// Send boardcast message to all member in room
		case message := <-room.broadcast:
			room.broadcastToClientsInRoom(message)
		}
	}
}

// Add client to new room and notify to all member
func (room *Room) registerClientInRoom(client *Client) {
	room.NotifyClientJoined(client)
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
