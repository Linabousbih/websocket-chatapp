package main

import (
	"chatapp/auth"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type room struct {

	// holds all current clients in the room
	clients map[*client]bool

	// join is a channel for all clients wishing to join the room
	join chan *client

	// leave is a channel for all clients wishing to leave the room
	leave chan *client

	// forward is a channel that holds incoming messages that should be forwarded to the other clients.

	forward chan []byte
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)
		case msg := <-r.forward:
			for client := range r.clients {
				client.receive <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

var rooms = make(map[string]*room)
var mu sync.Mutex

func getRoom(name string) *room {
	mu.Lock()
	defer mu.Unlock()
	if r, ok := rooms[name]; ok {
		return r
	}
	r := newRoom()
	rooms[name] = r
	go r.run()
	return r
}
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	roomName := req.URL.Query().Get("room")
	if roomName == "" {
		http.Error(w, "Room name required", http.StatusBadRequest)
		return
	}

	token := req.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	// ✅ Validate token and extract claims
	claims, err := auth.ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	username := claims.Email // or claims.UserID if you prefer

	realRoom := getRoom(roomName)

	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &client{
		socket:  socket,
		receive: make(chan []byte, messageBufferSize),
		room:    realRoom,
		name:    username,
	}

	realRoom.join <- client
	defer func() { realRoom.leave <- client }()
	go client.write()
	client.read()
}
