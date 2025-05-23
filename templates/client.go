package templates

import "github.com/gorilla/websocket"

// client represents a single chatting user
type client struct {

	// a socket is the web socket for this uer
	socket *websocket.Conn

	//receiver is a channel to receive messages from other clients
	receive chan []byte

	//room is the room this client is chatting in
	room *room
}

func (c *client) read() {
	// close the connection when we are done
	defer c.socket.Close()
	// endlessly read messages
	for {
		_, msg, err := c.socket.ReadMessage()
		// break if there is an error
		if err != nil {
			return
		}
		// forward the message to the room
		c.room.forward <- msg
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
