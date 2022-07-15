package types

import (
	nodesinterfaces "github.com/vitelabs/vite-portal/internal/nodes/interfaces"
)

// SessionHeader defines the header for session information
type SessionHeader struct {
	IpAddress   string `json:"ipAddress"`
	Chain       string `json:"chain"`
	RequestTime int64  `json:"requestTime"`
}

// NewSession creates a new session from seed data
func NewSession(keeper nodesinterfaces.KeeperI) {
	// TODO: move to nodes or session keeper?
}

// NewSessionNodes creates nodes for the session
func NewSessionNodes(keeper nodesinterfaces.KeeperI) {
	// TODO: move to nodes or session keeper?
}
