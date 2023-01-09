package ws

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func NewUpgrader(timeout time.Duration) *websocket.Upgrader {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// TODO: pass array of acceptable origins
			return true
		},
	}
	if timeout > 0 {
		upgrader.HandshakeTimeout = timeout
	}
	return &upgrader
}