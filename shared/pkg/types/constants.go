package types

import "time"

const (
	// The fallback IP address of the client (only used if configured TrueClientIp is not set)
	HeaderFallbackClientIp = "VP-Fallback-Client-IP"
	// The expiry timeout of JSON Web Tokens
	JWTExpiryTimeout = 60 * time.Second
	// The subject representing a relayer used in JSON Web Tokens
	JWTRelayerSubject = "vite-portal-relayer"
	// The maximum size of the payload in bytes
	MaxPayloadSize = 1024 * 128
	// Time allowed to read the next pong message from the peer.
	WebSocketPongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	WebSocketPingPeriod = (WebSocketPongWait * 9) / 10
)