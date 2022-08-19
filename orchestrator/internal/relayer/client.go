package relayer

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// Time allowed to write a message to the client.
	WriteWait time.Duration
	// The relayer hub.
	Hub *Hub
	// The websocket connection.
	Conn *websocket.Conn
	// Buffered channel of outbound messages.
	Send chan []byte
}

// ReadPump pumps messages from the websocket connection to the hub.
//
// The application runs ReadPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(types.MaxPayloadSize)
	c.Conn.SetReadDeadline(time.Now().Add(types.WebSocketPongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(types.WebSocketPongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// TODO: handle message
		logger.Logger().Info().Msg(fmt.Sprintf("%#v", message))
	}
}

// WritePump pumps messages from the hub to the websocket connection.
//
// A goroutine running WritePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(types.WebSocketPingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(c.WriteWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(c.WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
