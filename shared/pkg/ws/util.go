package ws

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func CanConnect(url string, timeout time.Duration) bool {
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: timeout,
	}
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func CloseConnection(c *websocket.Conn, msg string, timeout time.Duration) error {
	data := websocket.FormatCloseMessage(websocket.CloseNormalClosure, msg)
	err := c.WriteControl(websocket.CloseMessage, data, time.Now().Add(timeout))
	if err != nil && err != websocket.ErrCloseSent {
		// If close message could not be sent, then close without the handshake.
		return c.Close()
	}
	return nil
}
