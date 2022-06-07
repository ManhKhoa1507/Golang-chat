package main

import (
	"chat/config"
	"chat/models"
	"encoding/json"
	"fmt"
	"net/http"
)

const PubSubGeneralChannel = "general"

// Declare WsServer struct
type WsServer struct {
	clients        map[*Client]bool
	register       chan *Client
	unregister     chan *Client
	broadcast      chan []byte
	rooms          map[*Room]bool
	users          []models.User
	roomRepository models.RoomRepository
	userRepository models.UserRepository
}

// Create new WsServer type
func NewWebsocketServer(roomRepository models.RoomRepository, userRepository models.UserRepository) *WsServer {
	wsServer := &WsServer{
		clients:        make(map[*Client]bool),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		rooms:          make(map[*Room]bool),
		roomRepository: roomRepository,
		userRepository: userRepository,
	}

	// Add users from database to server
	wsServer.users = userRepository.GetAllUsers()

	return wsServer
}

// Server websocket
func ServerWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {
	// Get name from URL query
	name, ok := r.URL.Query()["name"]

	if !ok || len(name[0]) < 1 {
		print("Missing name")
		return
	}

	// Upgrade connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)

	// Error handle
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create new client and print result
	client := newClient(conn, wsServer, name[0])

	go client.writePump()
	go client.readPump()

	// register client
	wsServer.register <- client
}

// Redis pub/sub function
// Publish client joined the repo
func (server *WsServer) publishClientJoined(client *Client) {
	// Create Joined message
	message := &Message{
		Action: UserJoinedAction,
		Sender: client,
	}

	if err := config.Redis.Publish(ctx, PubSubGeneralChannel, message.encode()).Err(); err != nil {
		fmt.Println("Error when publish client joined")
	}
}

// Publish client left the repo
func (server *WsServer) publishClientLeft(client *Client) {
	// Create left message
	message := &Message{
		Action: UserLeftAction,
		Sender: client,
	}

	if err := config.Redis.Publish(ctx, PubSubGeneralChannel, message.encode()).Err(); err != nil {
		fmt.Println("Error when publish client left")
	}
}

// Listen to pub/sub general channel
func (server *WsServer) listenPubSubChannel() {
	pubsub := config.Redis.Subscribe(ctx, PubSubGeneralChannel)
	ch := pubsub.Channel()

	// Change to json format
	for msg := range ch {
		var message Message
		// Handle error
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			fmt.Println("Error when change message -> json format")
			return
		}

		// Actions
		switch message.Action {
		case UserJoinedAction:
			server.handleUserJoined(message)

		case UserLeftAction:
			server.handleUserLeft(message)

		case JoinRoomPrivateAction:
			server.handleUserJoinPrivate(message)
		}
	}
}
