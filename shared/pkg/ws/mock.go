package ws

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
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

func NewMockWsRpc(port int32) *MockWsRpc {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		panic(err)
	}

	var pattern = "/ws/mock"
	var realPort = l.Addr().(*net.TCPAddr).Port

	return &MockWsRpc{
		listener: l,
		Port: realPort,
		Pattern: pattern,
		Url: fmt.Sprintf("ws://localhost:%d%s", realPort, pattern),
	}
}

func (r *MockWsRpc) Serve(timeout time.Duration) error {
	hub := NewHub(timeout, messageHandler)
	go hub.Run()

	serveMux := http.NewServeMux()
	serveMux.HandleFunc(r.Pattern, func(w http.ResponseWriter, r *http.Request) {
		handleClient(hub, w, r, timeout)
	})

	return http.Serve(r.listener, serveMux)
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
