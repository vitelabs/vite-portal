package app

import (
	"time"

	"github.com/vitelabs/vite-portal/relayer/internal/orchestrator"
)

func InitOrchestrator(url string, timeout int64) *orchestrator.Orchestrator {
	return orchestrator.InitOrchestrator(url, time.Duration(timeout) * time.Millisecond)
} 