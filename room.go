package main

import (
	"chatapp/templates"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct{

	// holds all current clients in the room
	clients map[*client]bool

	// join is a channel for all clients wishing to join the room
	join chan *client

	// leave is a channel for all clients wishing to leave the room
	leave chan *client

	// forward is a channel that holds incoming messages that should be forwarded to the other clients.

	forward chan []byte 
}

func newRoom() *room{
	return &room{
		forward: make(chan []byte), 
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool)
	}
}

func (r * room) run(){
	for {
		select {
		case client := <-r.join:
			r.clients[client]=true
		case client := <- r.leave:
			delete(r.clients,client)
			close(client.receive)
		case msg :=<- r.forward:
			for client := range r.client{
				client.receive <-msg
			}
		}
	}
}

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServerHTTP(w http.ResponseWriter, req http.Request){
	socket,err:=upgrader.Upgrade(w,req)
	if err!=nil{
		log.Fatal("ServerHTTP:",err)
		return
	}

	client:=&client{
		socket: socket,
		receiver: make(chan[]byte, messageBufferSize),
		room: r
	}

	r.join<-client
	defer func(){r.leave <- client}()
	go client.write()
	client.read()
}