package types

import (
	nodesinterfaces "github.com/vitelabs/vite-portal/internal/nodes/interfaces"
	roottypes "github.com/vitelabs/vite-portal/internal/types"
)

// A read/write request to be relayed
type Relay struct {
	Payload Payload `json:"payload"` // The data payload of the request
}

// The data being sent to the external node
type Payload struct {
	Data    string            `json:"data"`              // The actual data string to be sent
	Method  string            `json:"method"`            // The HTTP CRUD method
	Path    string            `json:"path"`              // The REST path
	Headers map[string]string `json:"headers,omitempty"` // The HTTP headers
}

// The response to a relay request
type RelayResponse struct {
	Response string
}

// The response object used in dispatching
type DispatchResponse struct {
	Session DispatchSession `json:"session"`
}

type DispatchSession struct {
	Header SessionHeader `json:"header"`
	// Key SessionKey `json:"key"`
	Nodes []nodesinterfaces.NodeI `json:"nodes"`
}

// Execute attempts to do a request on the specified blockchain
func (r Relay) Execute() (string, roottypes.Error) {
	url := ""
	res, err := executeHttpRequest(r.Payload.Data, url, "", r.Payload.Method, r.Payload.Headers)
	if err != nil {
		// TODO: track metrics
		return res, NewError(ModuleName, CodeHttpExecutionError, err)
	}
	return res, nil
}

// executeHttpRequest takes in the raw json string and forwards it to the RPC endpoint
func executeHttpRequest(payload, url, userAgent string, method string, headers map[string]string) (string, error) {
	return "", nil
}
