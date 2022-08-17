package app

import (
	"github.com/vitelabs/vite-portal/relayer/internal/orchestrator/client"
	"github.com/vitelabs/vite-portal/relayer/internal/orchestrator/interfaces"
)

func NewOrchestrator() (interfaces.ClientI, error) {
	orchestrator := client.NewClient()
	return orchestrator, nil
}

func InitOrchestrator() (interfaces.ClientI, error) {
	orchestrator, err := NewOrchestrator()
	if err != nil {
		return nil, err
	}
	err = orchestrator.Connect()
	if err != nil {
		return nil, err
	}
	return orchestrator, nil
}