package chat

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a single chatting user.
type Client struct {
	Conn *websocket.Conn
	Pool *Pool
}

// Pool represents a pool of WebSocket connections.
type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
	sync.Mutex
}

// Message represents a message sent by a client.
type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

// NewPool creates a new WebSocket pool.
func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

// Start initializes the pool and starts listening for register, unregister, and broadcast events.
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
		case message := <-pool.Broadcast:
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					return
				}
			}
		}
	}
}

// Read listens for incoming messages from the WebSocket connection.
func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}
		msg := Message{Type: 1, Body: string(message)}
		c.Pool.Broadcast <- msg
	}
}

// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
