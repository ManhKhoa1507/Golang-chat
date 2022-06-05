package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Define user action
	UserJoinedAction = "user-join"
	UserLeftAction   = "user-left"
)

// Create upgrader to hold the buffer size of websocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

// define new line and space character
var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// represents the websocket client at server
type Client struct {
	// The actual websocket connection.
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	rooms    map[*Room]bool
}

// Define some variables
const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

// Run client
func (server *WsServer) Run() {
	// Create endless loop to listen
	for {
		select {

		// Register client
		case client := <-server.register:
			server.registerClient(client)

		// Unregister client
		case client := <-server.unregister:
			server.unregisterClient(client)

		// Broadcast message
		case message := <-server.broadcast:
			server.broadcastToClients(message)
		}
	}
}

// Client interactive
// return new websocket connection
func newClient(conn *websocket.Conn, wsServer *WsServer, name string) *Client {
	fmt.Println("Client " + name + " joined")
	return &Client{
		ID:       uuid.New(),
		Name:     name,
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		rooms:    make(map[*Room]bool),
	}
}

// Function ro register client -> enable flag true to server.clients[client]
func (server *WsServer) registerClient(client *Client) {
	server.notifyClientJoined(client)
	server.listOnlineClients(client)
	server.clients[client] = true
}

// Function to unregisterClient -> delete user
func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
		server.notifyClientLeft(client)
	}
}

// Disconnect client
func (client *Client) disconnect() {

	// add client to unregister and close connection
	client.wsServer.unregister <- client

	// Unregister client from room
	for room := range client.rooms {
		room.unregister <- client
	}

	close(client.send)
	client.conn.Close()
}

// Boardcast message to Clients
func (server *WsServer) broadcastToClients(message []byte) {
	fmt.Println("Message: ", message)
	// Send message to all client in server.clients
	for client := range server.clients {
		client.send <- message
	}
}

// Get client information
// Get client name
func (client *Client) getName() string {
	return client.Name
}

// Get client ID
func (client *Client) getID() string {
	return client.ID.String()
}

// Find client by ID
func (server *WsServer) findClientByID(ID string) *Client {
	var foundClient *Client
	for client := range server.clients {
		if client.getID() == ID {
			foundClient = client
			break
		}
	}
	return foundClient
}

// Notify memeber joined & left
// Notify new client joined
func (server *WsServer) notifyClientJoined(client *Client) {
	// Create welcome new member message
	message := &Message{
		Action: UserJoinedAction,
		Sender: client,
	}

	// Broadcast to all client
	server.broadcastToClients(message.encode())
}

// Notify client left the room
func (server *WsServer) notifyClientLeft(client *Client) {
	// Create goodbye message
	message := &Message{
		Action: UserLeftAction,
		Sender: client,
	}

	// Broadcast to all client
	server.broadcastToClients(message.encode())
}

// ReadPump and WritePump
// Read message and send over the Websocket connection, make endless loop untils client is disconnected
func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	// Check for connection
	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless loop to waiting message
	for {
		// Read message
		_, jsonMessage, err := client.conn.ReadMessage()

		// Error Handle
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		// add jsonMessage to broadcast
		client.handleNewMessage(jsonMessage)
	}
}

// Goroutine handles sending messages to the connected client
// Run endless loop waiting for new message in client.send channel
// When received new message -> writes to clients
// If multiple message, combined in one write
func (client *Client) writePump() {
	// Use NewTicker to viewing time
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	// Create endless loop
	for {
		select {

		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(pongWait))

			// If WebSocket closed the channel
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// NextWriter return a writer for next message to send
			// return writer w
			w, err := client.conn.NextWriter(websocket.TextMessage)

			// Handle the error
			if err != nil {
				fmt.Println("Error")
				return
			}

			// Writer write message
			w.Write(message)

			// Attach queued chat message to the current websocket message
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			// Handle the error when writer closed
			if err := w.Close(); err != nil {
				return
			}

		// Handle connection
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
