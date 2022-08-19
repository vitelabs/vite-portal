package rpc

import (
	"net/http"
	"time"

	"github.com/vitelabs/vite-portal/shared/pkg/util/httputil"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
)

func handleRelayer(hub *ws.Hub, w http.ResponseWriter, r *http.Request, timeout int64) {
	// TODO: add authorization
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
