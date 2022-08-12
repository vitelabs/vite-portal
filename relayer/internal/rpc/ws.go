package rpc

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/vitelabs/vite-portal/internal/logger"
	"github.com/vitelabs/vite-portal/internal/types"
	"github.com/zyedidia/generic/mapset"
	"golang.org/x/net/websocket"
)

func StartWsRpc(port int32, timeout int64) {
	server := NewServer()
	handler := server.WebsocketHandler([]string{"*"})

	serveMux := http.NewServeMux()
	serveMux.Handle("/ws/v1", handler)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), serveMux)
	if err != nil {
		logger.Logger().Error().Err(err).Msg("WS RPC closed")
		os.Exit(1)
	}
}

func (s *Server) WebsocketHandler(allowedOrigins []string) http.Handler {
	return websocket.Server{
		Handshake: handshakeValidator(allowedOrigins),
		Handler: func(conn *websocket.Conn) {
			// Create custom encode/decode pair to enforce payload size
			conn.MaxPayloadBytes = types.MaxRequestContentLength

			encoder := func(v interface{}) error {
				return defaultCodec.Send(conn, v)
			}
			decoder := func(v interface{}) error {
				return defaultCodec.Receive(conn, v)
			}
			
			//safeConn := deadliner{conn, time.Millisecond*100}

			s.ServeCodec(NewCodec(conn, encoder, decoder))
		},
	}
}

// handshakeValidator returns a handler that verifies the origin during the
// websocket upgrade process. If '*' is specified as an allowed origin all
// connections are accepted.
func handshakeValidator(allowedOrigins []string) func(*websocket.Config, *http.Request) error {
	origins := mapset.New[string]()
	allowAllOrigins := false

	for _, origin := range allowedOrigins {
		if origin == "*" {
			allowAllOrigins = true
		}
		if origin != "" {
			origins.Put(strings.ToLower(origin))
		}
	}

	// allow localhost if no allowedOrigins are specified.
	if origins.Size() == 0 {
		origins.Put("http://localhost")
		if hostname, err := os.Hostname(); err == nil {
			origins.Put("http://" + strings.ToLower(hostname))
		}
	}

	logger.Logger().Debug().Str("origins", fmt.Sprintf("%#v", origins)).Msg("Allowed origin(s) for WS RPC interface.")

	f := func(cfg *websocket.Config, req *http.Request) error {
		if req.ContentLength > types.MaxRequestContentLength {
			return fmt.Errorf("content exceeds the limit of %d bytes", types.MaxRequestContentLength)
		}

		origin := strings.ToLower(req.Header.Get("Origin"))
		if allowAllOrigins || origins.Has(origin) {
			return nil
		}
		logger.Logger().Warn().Msg(fmt.Sprintf("origin '%s' not allowed on WS-RPC interface\n", origin))
		return fmt.Errorf("origin %s not allowed", origin)
	}

	return f
}

var defaultCodec = websocket.Codec{
	Marshal: func(v interface{}) ([]byte, byte, error) {
		switch data := v.(type) {
		case string:
			return []byte(data), websocket.TextFrame, nil
		case []byte:
			return data, websocket.BinaryFrame, nil
		}
		return nil, websocket.UnknownFrame, websocket.ErrNotSupported
	},
	Unmarshal: func(msg []byte, payloadType byte, v interface{}) error {
		switch data := v.(type) {
		case *string:
			*data = string(msg)
			return nil
		case *[]byte:
			*data = msg
			return nil
		}
		return websocket.ErrNotSupported
	},
}

// deadliner is a wrapper around net.Conn that sets read/write deadlines before
// every Read() or Write() call.
// Source: https://github.com/gobwas/ws-examples/blob/master/src/chat/main.go
type deadliner struct {
	net.Conn
	t time.Duration
}

func (d deadliner) Write(p []byte) (int, error) {
	if err := d.Conn.SetWriteDeadline(time.Now().Add(d.t)); err != nil {
		return 0, err
	}
	return d.Conn.Write(p)
}

func (d deadliner) Read(p []byte) (int, error) {
	if err := d.Conn.SetReadDeadline(time.Now().Add(d.t)); err != nil {
		return 0, err
	}
	return d.Conn.Read(p)
}