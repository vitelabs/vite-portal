package ws

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/httputil"
)

type MockWsRpc struct {
	listener net.Listener
	Port int
	Pattern string
	Url string
}

func (r *MockWsRpc) Close() error {
	return r.listener.Close()
}

func StartMockWsRpc(timeout time.Duration) *MockWsRpc {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	var pattern = "/ws/mock"
	var port = l.Addr().(*net.TCPAddr).Port

	go startMockWsRpc(l, pattern, timeout)

	return &MockWsRpc{
		listener: l,
		Port: port,
		Pattern: pattern,
		Url: fmt.Sprintf("ws://localhost:%d%s", port, pattern),
	}
}

func startMockWsRpc(l net.Listener, pattern string, timeout time.Duration) {
	hub := NewHub(timeout, messageHandler)
	go hub.Run()

	serveMux := http.NewServeMux()
	serveMux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		handleClient(hub, w, r, timeout)
	})

	err := http.Serve(l, serveMux)
	if err != nil {
		panic(err)
	}
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
