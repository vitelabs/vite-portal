package service

import (
	nodeinterfaces "github.com/vitelabs/vite-portal/internal/node/interfaces"
)

// Service maintains the link to storage and exposes getter/setter methods related to core functionalities
type Service struct {
	NodeService nodeinterfaces.ServiceI
}

// NewService creates new instances of the core module service
func NewService(nodeService nodeinterfaces.ServiceI) Service {
	return Service{
		NodeService: nodeService,
	}
}
