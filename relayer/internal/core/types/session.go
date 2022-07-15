package types

import (
	nodeinterfaces "github.com/vitelabs/vite-portal/internal/node/interfaces"
)

// SessionHeader defines the header for session information
type SessionHeader struct {
	IpAddress   string `json:"ipAddress"`
	Chain       string `json:"chain"`
	RequestTime int64  `json:"requestTime"`
}

// NewSession creates a new session from seed data
func NewSession(s nodeinterfaces.ServiceI) {
	// TODO: move to nodes or session service?
}

// NewSessionNodes creates nodes for the session
func NewSessionNodes(s nodeinterfaces.ServiceI) {
	// TODO: move to nodes or session service?
}
