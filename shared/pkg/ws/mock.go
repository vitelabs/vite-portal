package ws

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/httputil"
)

type MockWsRpc struct {
	cancel context.CancelFunc
	listener net.Listener
	Port int
	Pattern string
	Url string
}

func (r *MockWsRpc) Close() error {
	r.cancel()
	return r.listener.Close()
}

func StartMockWsRpc(port int32, timeout time.Duration) *MockWsRpc {
	ctx, cancel := context.WithCancel(context.Background())

	l, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		panic(err)
	}

	var pattern = "/ws/mock"
	var realPort = l.Addr().(*net.TCPAddr).Port

	go startMockWsRpc(ctx, l, pattern, timeout)

	return &MockWsRpc{
		cancel: cancel,
		listener: l,
		Port: realPort,
		Pattern: pattern,
		Url: fmt.Sprintf("ws://localhost:%d%s", realPort, pattern),
	}
}

func startMockWsRpc(ctx context.Context, l net.Listener, pattern string, timeout time.Duration) {
	hub := NewHub(timeout, messageHandler)
	go hub.Run()

	serveMux := http.NewServeMux()
	serveMux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		handleClient(hub, w, r, timeout)
	})

	http.Serve(l, serveMux)
}

func messageHandler(client *Client, msg []byte) {
	logger.Logger().Info().Msg(fmt.Sprintf("%s", msg))
	client.Send <- msg
}

func handleClient(hub *Hub, w http.ResponseWriter, r *http.Request, timeout time.Duration) {
	err := hub.RegisterClient(w, r, timeout)
	if err != nil {
		logger.Logger().Error().Err(err).Msg("register client failed")
		httputil.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
