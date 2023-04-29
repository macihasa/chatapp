package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WEBSOCKET HANDLERS
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // Allow all origins
}

// The distibution hub keeps track of all
// active clients and broadcasts messages between them
type DistributionHub struct {
	clients    map[client]bool
	broadcast  chan []byte
	register   chan *client
	unregister chan *client
}

// Clients are registered to the distribution hub. When a message is received from the client connection
// it gets pushed to the broadcast channel of the distibution hub and sent to all other active clients
type client struct {
	conn *websocket.Conn
	send chan []byte
}

// Initialize hub
var hub = DistributionHub{
	clients:    make(map[client]bool),
	broadcast:  make(chan []byte),
	register:   make(chan *client),
	unregister: make(chan *client),
}

func StartDistributionHub() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[*client] = true

		case client := <-hub.unregister:
			if _, ok := hub.clients[*client]; ok {
				delete(hub.clients, *client)
				close(client.send)
			}
		case message := <-hub.broadcast:
			for client := range hub.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(hub.clients, client)
				}
			}
		}
	}
}

// WsNewClient creates a new client and registeres it to the DistributionHub
func WsNewClient(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error: ", err)
		fmt.Fprint(w, err)
		return
	}

	client := &client{
		conn: conn,
		send: make(chan []byte),
	}

	hub.register <- client
	log.Printf("\n-----Registered client: %v-----\n", client)
	defer func() {
		log.Printf("-----Unregistered client: %v-----\n", client)
		hub.unregister <- client
	}()

	// start goroutines that listens for incoming and outgoing messages
	go sendMessage(client)
	receiveMessage(client, w)
}

func sendMessage(client *client) {
	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				hub.unregister <- client
				err := client.conn.Close()
				if err != nil {
					log.Println(err)
					break
				}
			}
			client.conn.WriteMessage(1, message)
		}
	}
}

func receiveMessage(client *client, w http.ResponseWriter) {
	defer func() {
		hub.unregister <- client
		client.conn.Close()
	}()
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			fmt.Fprintln(w, err)
			break
		}
		fmt.Println(string(message))
		hub.broadcast <- message
	}
}
