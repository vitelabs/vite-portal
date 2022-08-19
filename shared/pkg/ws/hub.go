package ws

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
}

func NewHub(msgHandler func(client *Client, msg []byte)) *Hub {
	return &Hub{
		Clients:        make(map[*Client]bool),
		Broadcast:      make(chan []byte),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		MessageHandler: msgHandler,
	}
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
