package main

import (
	"fmt"
	"log"
	"net/http"
	"rtchat/internal/chat"
)

func serveWs(pool *chat.Pool, w http.ResponseWriter, r *http.Request) {
	conn, err := chat.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
		return
	}

	client := &chat.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := chat.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
}

func main() {
	fmt.Println("Real-Time Chat App v0.01")
	setupRoutes()
	log.Fatal(http.ListenAndServe("127.0.0.1:9090", nil))
}
