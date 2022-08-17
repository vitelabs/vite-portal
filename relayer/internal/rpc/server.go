package rpc

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
)

type Server struct {
	run      int32
	codecsMu sync.Mutex
	codecs   map[string]ServerCodec
}

func NewServer() *Server {
	server := &Server{
		codecs: map[string]ServerCodec{},
		run:    1,
	}

	return server
}

// ServeCodec reads incoming requests from codec, calls the appropriate callback and writes the
// response back using the given codec. It will block until the codec is closed or the server is
// stopped. In either case the codec is closed.
func (s *Server) ServeCodec(codec ServerCodec) error {
	defer codec.Close()
	return s.serveRequest(context.Background(), codec)
}

// serveRequest reads requests from the codec, calls the callback and
// writes the response to the given codec.
//
// Handles requests until the codec returns an error when reading (in most cases
// an EOF). It executes requests in parallel.
func (s *Server) serveRequest(ctx context.Context, codec ServerCodec) error {
	var pend sync.WaitGroup

	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			logger.Logger().Error().Str("error", string(buf)).Msg("recover error")
		}
		s.codecsMu.Lock()
		delete(s.codecs, codec.Id())
		s.codecsMu.Unlock()
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s.codecsMu.Lock()
	if atomic.LoadInt32(&s.run) != 1 { // server stopped
		s.codecsMu.Unlock()
		return &shutdownError{}
	}
	s.codecs[codec.Id()] = codec
	s.codecsMu.Unlock()

	// test if the server is ordered to stop
	for atomic.LoadInt32(&s.run) == 1 {
		reqs, err := codec.ReadRequest()
		if err != nil {
			// If read error occurred, send an error
			if err.Error() != "EOF" {
				logger.Logger().Error().Err(err).Msg("read error")
				codec.Write(codec.CreateErrorResponse(err))
			}
			// Error or end of stream, wait for requests and tear down
			pend.Wait()
			return nil
		}

		// check if server is ordered to shutdown and return an error
		// telling the client that his request failed.
		if atomic.LoadInt32(&s.run) != 1 {
			err = &shutdownError{}
			codec.Write(codec.CreateErrorResponse(err))
			return nil
		}

		// Start goroutine to serve and loop back
		pend.Add(1)

		go func(reqs []byte) {
			defer pend.Done()
			s.exec(ctx, codec, reqs)
		}(reqs)
	}
	return nil
}

// exec executes the given request and writes the result back using the codec.
func (s *Server) exec(ctx context.Context, codec ServerCodec, req []byte) {
	if err := codec.Write("asdf"); err != nil {
		logger.Logger().Error().Err(err).Msg("server write error")
		codec.Close()
	}
}

// Error wraps RPC errors, which contain an error code in addition to the message.
type Error interface {
	Error() string  // returns the message
	ErrorCode() int // returns the code
}

type ServerCodec interface {
	// The identifier of the codec
	Id() string
	// Read next request
	ReadRequest() ([]byte, Error)
	// Assemble error response
	CreateErrorResponse(err Error) interface{}
	// Write msg to client
	Write(msg interface{}) error
	// Close underlying data stream
	Close()
	// Closed when underlying connection is closed
	Closed() <-chan interface{}
}
