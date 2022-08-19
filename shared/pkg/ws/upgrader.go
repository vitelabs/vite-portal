package ws

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func NewUpgrader(timeout int64) *websocket.Upgrader {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	if timeout > 0 {
		upgrader.HandshakeTimeout = time.Duration(timeout) * time.Millisecond
	}
	return &upgrader
}