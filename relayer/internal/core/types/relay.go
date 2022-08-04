package types

// A read/write request to be relayed
type Relay struct {
	Host     string  `json:"host"`
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

// The response of a single node
type NodeResponse struct {
	NodeId       string `json:"nodeId"`
	ResponseTime int64  `json:"responseTime"`
	Response     string `json:"response"`
	Error        string  `json:"error"`
}

// The result of a relay request
type RelayResult struct {
	SessionKey string         `json:"sessionKey"`
	Relay      Relay          `json:"relay"`
	Responses  []NodeResponse `json:"responses"`
}

// The response to a relay request
type RelayResponse struct {
	SessionKey string `json:"sessionKey"`
	Response   string `json:"response"`
}
