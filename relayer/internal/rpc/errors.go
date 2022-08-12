package rpc

// issued when a request is received after the server is issued to stop.
type shutdownError struct{}

func (e *shutdownError) ErrorCode() int { return -32000 }

func (e *shutdownError) Error() string { return "server is shutting down" }

// received message isn't a valid request
type invalidRequestError struct {
	message string
}

func (e *invalidRequestError) ErrorCode() int { return -32600 }

func (e *invalidRequestError) Error() string { return e.message }