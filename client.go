package main

// TODO: Continue from page 20 when connected to network
// go-programming-blueprints-2nd.pdf
// TODO: go get github.com/gorilla/websocket
import (
	"time"

	"github.com/gorilla/websocket"
)

// client object
// socket 	gets reference to websocket object
// send 	creates an array of bytes
// room 	get reference to room object
type client struct {
	// socket is the websocket for this client
	socket *websocket.Conn
	// send is a channel on which messages are sent
	send chan *message
	// room is the room this client is chating in
	room *room
	// userData holds the infomation about the users
	userData map[string]interface{}
}

// read method
// the client socket is closed until the surrounding function returns
// msg is the second returned value from the stablished client's socket
// the client byte array receives the message from the socket
func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		if avatarURL, ok := c.userData["avatar_url"]; ok {
			msg.AvatarURL = avatarURL.(string)
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
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}
