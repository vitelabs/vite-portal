package ws

type ConnectionStatus int64

const (
	Unknown ConnectionStatus = iota
	Connecting
	Connected
	Disconnected
)