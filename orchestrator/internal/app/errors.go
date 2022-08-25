package app

import "errors"

var (
	ErrDatadirUsed         = errors.New("datadir already used by another process")
	ErrOrchestratorStopped = errors.New("orchestrator not started")
	ErrOrchestratorRunning = errors.New("orchestrator already running")
	ErrServiceUnknown      = errors.New("unknown service")
)
