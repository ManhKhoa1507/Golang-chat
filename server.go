package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("localhost", ":8080", "http server address")

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
