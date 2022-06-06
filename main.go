package main

import (
	"chat/config"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	flag.Parse()

	// upgrade connection to websocket
	wsServer := NewWebsocketServer()
	go wsServer.Run()

	// Handle request to this part
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServerWs(wsServer, w, r)
	})

	// Host a fileserver dir from ./public
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// Open connection
	fmt.Println("Server running at: localhost", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))

	// Create database
	db := config.InitDB()
	defer db.Close()
}