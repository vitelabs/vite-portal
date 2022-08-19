package rpc

import (
	"fmt"
	"net/http"
	"os"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/httputil"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
)

func StartWsRpc(port int32, timeout int64) {
	hub := ws.NewHub(timeout, messageHandler)
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
	err := hub.RegisterClient(w, r, timeout)
	if err != nil {
		logger.Logger().Error().Err(err).Msg("register client failed")
		httputil.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
