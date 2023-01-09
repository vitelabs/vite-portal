package ws

// issued when a request is received after the server is issued to stop.
type ShutdownError struct{}

func (e *ShutdownError) ErrorCode() int { return -32000 }

func (e *ShutdownError) Error() string { return "server is shutting down" }

// received message isn't a valid request
type InvalidRequestError struct {
	message string
}

func (e *InvalidRequestError) ErrorCode() int { return -32600 }

func (e *InvalidRequestError) Error() string { return e.message }