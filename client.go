package main

// TODO: Continue from page 20 when connected to network
// go-programming-blueprints-2nd.pdf
// TODO: go get github.com/gorilla/websocket
import (
	"github.com/gorilla/websocket"
)

// client object
// socket 	gets reference to websocket object
// send 	creates an array of bytes
// room 	get reference to room object
type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

// read method
// the client socket is closed until the surrounding function returns
// msg is the second returned value from the stablished client's socket
// the client byte array receives the message from the socket
func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

// write method
// the client socket is closed until the surrounding function returns
// the loop reads every msg from the byte stream
// the WriteMessage from the socket returns an error in case the msg is not correctly conveted to TextMessage
// the error will be nil in case it was successful
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
