package types

import "time"

const (
	// The maximum size of the payload in bytes
	MaxPayloadSize = 1024 * 128
	// Time allowed to read the next pong message from the peer.
	WebSocketPongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	WebSocketPingPeriod = (WebSocketPongWait * 9) / 10
)