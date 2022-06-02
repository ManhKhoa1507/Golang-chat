package main

import (
	"encoding/json"
	"fmt"
)

// Define Action
const (
	// Define Server action
	SendMessageAction = "send-message"
	JoinRoomAction    = "join-room"
	LeaveRoomAction   = "leave-room"

	// Define user action
	UserJoinedAction = "user-join"
	UserLeftAction   = "user-left"
)

// Define message json
type Message struct {
	Action  string  `json:"action"`
	Message string  `json:"message"`
	Target  string  `json:"target"`
	Sender  *Client `json:"sender"`
}

// Encode message to json format
func (message *Message) encode() []byte {
	json, err := json.Marshal(message)

	// Handle error
	if err != nil {
		fmt.Println("Error")
	}
	return json
}

// Handle new message from client
func (client *Client) handleNewMessage(jsonMessage []byte) {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		fmt.Println("Error on unmarshal message -> JSON")
	}

	// Assign sender
	message.Sender = client

	switch message.Action {

	// If action is send message -> send message to specific room
	case SendMessageAction:
		roomName := message.Target

		// Find room to send
		if room := client.wsServer.findRoomByName(roomName); room != nil {
			room.broadcast <- &message
		}

	// If client requests to join new room
	case JoinRoomAction:
		client.handleJoinRoomMessage(message)

	// If client request to leaves room
	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message)
	}
}
