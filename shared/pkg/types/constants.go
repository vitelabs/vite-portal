package types

import "time"

const (
	// The fallback IP address of the client (only used if configured TrueClientIp is not set)
	HTTPHeaderFallbackClientIp = "VP-Fallback-Client-IP"
	// Can be used to provide credentials that authenticate a user agent with a server, allowing access to a protected resource.
	HTTPHeaderAuthorization = "Authorization"
	// The issuer representing a relayer used in JSON Web Tokens
	JWTRelayerIssuer = "vite-portal-relayer"
	// The default JWT secret
	DefaultJwtSecret = "secret1234"
	// The default JWT expiry timeout
	DefaultJwtExpiryTimeout = 0
	// The default Kafka server
	DefaultKafkaServer = "localhost:9092"
	// The Kafka group identifier related to online status of fullnodes
	DefaultKafkaGroupId = "vote_fullnode_award_event_prod"
	// The Kafka topic related to online status of fullnodes
	DefaultKafkaTopic = "fullnode_award_event_prod"
	// The Kafka group identifier related to RPC of fullnodes
	RpcKafkaGroupId = "vote_fullnode_rpc_event_prod"
	// The Kafka topic related to RPC of fullnodes
	RpcKafkaTopic = "fullnode_rpc_event_prod"
	// The maximum size of the payload in bytes
	MaxPayloadSize = 1024 * 128
	// Time allowed to read the next pong message from the peer.
	WebSocketPongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	WebSocketPingPeriod = (WebSocketPongWait * 9) / 10
)
