package types

// A read/write request to be relayed
type Relay struct {
	Chain    string  `json:"chain"`
	ClientIp string  `json:"clientIp"`
	Payload  Payload `json:"payload"`
}

// The data being sent to the external node
type Payload struct {
	Data    string              `json:"data"`              // The actual data string to be sent
	Method  string              `json:"method"`            // The HTTP CRUD method
	Path    string              `json:"path"`              // The REST path
	Headers map[string][]string `json:"headers,omitempty"` // The HTTP headers
}

// The response to a relay request
type RelayResponse struct {
	SessionKey string `json:"sessionKey"`
	Response   string `json:"response"`
}
