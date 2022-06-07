package main

import (
	"chat/config"
	"chat/repository"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
)

// Handle cancle signal for go routine
var (
	ctx  = context.Background()
	addr = flag.String("addr", ":8080", "http server addres")
)

func main() {
	flag.Parse()
	// Create Redis
	config.CreateRedisClient()

	// Create database
	db := config.InitDB()
	defer db.Close()

	// upgrade connection to websocket
	wsServer := NewWebsocketServer(&repository.RoomRepository{Db: db}, &repository.UserRepository{Db: db})
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
}
