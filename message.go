package main

import (
	"chat/models"
	"encoding/json"
	"fmt"
)

// Define Action
const (
	// Define Server action
	SendMessageAction = "send-message"
	JoinRoomAction    = "join-room"
	LeaveRoomAction   = "leave-room"

	// Define joined action
	RoomJoinedAction      = "room-joined"
	JoinRoomPrivateAction = "join-room-private"
)

// Define message json
type Message struct {
	Action  string      `json:"action"`
	Message string      `json:"message"`
	Target  *Room       `json:"target"`
	Sender  models.User `json:"sender"`
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

// UnmarshalJSON to create client instance for sender
func (message *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	msg := &struct {
		Sender Client `json:"sender"`
		*Alias
	}{
		Alias: (*Alias)(message),
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	message.Sender = &msg.Sender
	return nil
}

// Handle new message from client
func (client *Client) handleNewMessage(jsonMessage []byte) {
	var message Message

	// Handle error
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		fmt.Println("Error on unmarshal message -> JSON")
	}

	// Assign sender
	message.Sender = client

	switch message.Action {

	// If action is send message -> send message to specific room
	case SendMessageAction:
		roomID := message.Target.GetID()

		// Find room to send
		if room := client.wsServer.findRoomByID(roomID); room != nil {
			room.broadcast <- &message
		}

	// If client requests to join new room
	case JoinRoomAction:
		client.handleJoinRoomMessage(message)

	// If client request to leaves room
	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message)

	// If client request joined Private room
	case JoinRoomPrivateAction:
		client.handleJoinRoomPrivateMessage(message)
	}
}
