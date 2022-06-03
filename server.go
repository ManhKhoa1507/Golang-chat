package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

// Declare WsServer struct
type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	rooms      map[*Room]bool
}

// Create new WsServer type
func NewWebsocketServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		rooms:      make(map[*Room]bool),
	}
}

var addr = flag.String("localhost", ":8080", "http server address")

// Server websocket
func Serverws(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {
	// Upgrade connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)

	// Error handle
	if err != nil {
		fmt.Println(err)
		return
	}
	// Get name from URL query
	name, ok := r.URL.Query()["name"]
	if !ok || len(name[0]) < 1 {
		print("Missing name")
	}
	// Create new client and print result
	client := newClient(conn, wsServer, name[0])

	go client.writePump()
	go client.readPump()

	// register client
	wsServer.register <- client
}

func main() {
	flag.Parse()

	// upgrade connection to websocket
	wsServer := NewWebsocketServer()
	go wsServer.Run()

	// Handle request to this part
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		Serverws(wsServer, w, r)
	})

	// Host a fileserver dir from ./public
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// Open connection
	log.Fatal(http.ListenAndServe(*addr, nil))
}
