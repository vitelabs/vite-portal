package rpc

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/vitelabs/vite-portal/orchestrator/internal/app"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/httputil"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
)

func handleNode(hub *ws.Hub, w http.ResponseWriter, r *http.Request, timeout time.Duration) {
	// TODO: validate node
	chain, err := getChain(r)
	if err != nil {
		httputil.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Println(chain)

	err = hub.RegisterClient(w, r, timeout)
	if err != nil {
		logger.Logger().Error().Err(err).Msg("register client failed")
		httputil.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func getChain(r *http.Request) (string, error) {
	chain := r.URL.Query().Get("chain")
	if chain == "" {
		return app.CoreApp.Config.DefaultChain, nil
	}
	for _, v := range app.CoreApp.Config.SupportedChains {
		if chain == v {
			return chain, nil
		}
	}
	return "", errors.New("chain not supported")
}