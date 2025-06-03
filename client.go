package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

// client represents a single chatting user
type client struct {

	// a socket is the web socket for this uer
	socket *websocket.Conn

	//receive is a channel to receive messages from other clients
	receive chan []byte

	//room is the room this client is chatting in
	room *room

	name string
}

// Used to send messages
func (c *client) read() {
	// close the connection when we are done
	defer c.socket.Close()
	// endlessly read messages from input
	for {
		_, msg, err := c.socket.ReadMessage()
		// break if there is an error
		if err != nil {
			return
		}

		outgoing := map[string]string{
			"name":    c.name,
			"message": string(msg),
		}

		jsMessage, err := json.Marshal(outgoing)
		if err != nil {
			fmt.Println("Enconding failed!")
			continue
		}

		// forward the message to the room
		c.room.forward <- jsMessage
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.receive {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)

		if err != nil {
			return
		}
	}
}
