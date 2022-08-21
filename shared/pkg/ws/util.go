package ws

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func CanConnect(url string, timeout time.Duration) bool {
	dialer := websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
		HandshakeTimeout: timeout,
	}
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}