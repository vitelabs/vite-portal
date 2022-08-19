package rpc

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/httputil"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartWsRpc(port int32, timeout int64) {
	hub := ws.NewHub(messageHandler)
	go hub.Run()

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ws/v1", func(w http.ResponseWriter, r *http.Request) {
		handleClient(hub, w, r, timeout)
})

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), serveMux)
	if err != nil {
		logger.Logger().Error().Err(err).Msg("WS RPC closed")
		os.Exit(1)
	}
}

func messageHandler(client *ws.Client, msg []byte) {
	logger.Logger().Info().Msg(fmt.Sprintf("%s", msg))
	client.Send <- msg
}

func handleClient(hub *ws.Hub, w http.ResponseWriter, r *http.Request, timeout int64) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		httputil.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	client := &ws.Client{
		WriteWait: time.Duration(timeout) * time.Millisecond,
		Hub:       hub,
		Conn:      c,
		Send:      make(chan []byte, 256),
	}
	hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go client.WritePump()
	go client.ReadPump()
}
