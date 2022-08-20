package ws

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Hub maintains the set of active relayers and broadcasts messages to them.
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool
	// Inbound messages from the clients.
	Broadcast chan []byte
	// Register requests from the clients.
	Register chan *Client
	// Unregister requests from clients.
	Unregister chan *Client
	// The function handling incoming messages from clients.
	MessageHandler func(client *Client, msg []byte)
	// The HTTP connection upgrader
	upgrader *websocket.Upgrader
}

func NewHub(timeout time.Duration, msgHandler func(client *Client, msg []byte)) *Hub {
	hub := Hub{
		Clients:        make(map[*Client]bool),
		Broadcast:      make(chan []byte),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		MessageHandler: msgHandler,
		upgrader:       NewUpgrader(timeout),
	}
	return &hub
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

func (h *Hub) RegisterClient(w http.ResponseWriter, r *http.Request, timeout time.Duration) error {
	c, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	client := &Client{
		WriteWait: timeout,
		Hub:       h,
		Conn:      c,
		Send:      make(chan []byte),
	}
	h.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go client.WritePump()
	go client.ReadPump()

	return nil
}
