package main

import (
	"log"
	"net/http"

	"chatAppGo/trace"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

// chat room
// forward is a byte array channel
// join is a channel using a client reference
// leave is a channel using a client reference
// clients holds a boolean map that uses client references as keys
// tracer holds a Tracer object to report on tests
type room struct {
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	tracer  trace.Tracer
}

// newRoom is a helper function that returns the reference
// to a room
func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

// run is a function that waits on clients to join or leave
// or a byte array channel dequeues to a msg
func (r *room) run() {
	for {
		select {
		// the client receives a dequeued client that just joined the room
		case client := <-r.join:
			// the clients map uses the client reference as a key and set to true
			r.clients[client] = true
			r.tracer.Trace("New client joined")
		// the client receives a dequeued client that just left the room
		case client := <-r.leave:
			// client reference is deleted from the clients map
			delete(r.clients, client)
			// the byte array channel is closed
			close(client.send)
			r.tracer.Trace("Client left")
		// msg receives a byte array dequeued from the forward channel
		case msg := <-r.forward:
			r.tracer.Trace("Message received: ", msg.Message)
			// for every client in the client map, place the byte array in the send channel
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace(" -- sent to client")
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// reference to a websocket.Upgrader is given to upgrader
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

// ServeHTTP method using a room reference
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}
	// socket is connected to client
	// send is created as a channel array
	// client holds a room reference
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	// the reference to client is given to a room channel
	r.join <- client
	// the reference to client is not put in the leave channel until
	defer func() { r.leave <- client }()
	// a goroutine is started and called using write
	go client.write()
	client.read()
}
