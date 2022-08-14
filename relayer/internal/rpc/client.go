package rpc

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/vitelabs/vite-portal/internal/app"
	coretypes "github.com/vitelabs/vite-portal/internal/core/types"
	"github.com/vitelabs/vite-portal/internal/logger"
	"github.com/vitelabs/vite-portal/internal/types"
)

func Relay(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if cors(&w, r) {
		return
	}
	if logger.DebugEnabled() {
		logger.Logger().Debug().Str("request", fmt.Sprintf("%#v", r)).Msg("relay request")
	}
	relay, err1 := extractRelay(w, r, p)
	if err1 != nil {
		response := types.RpcRelayErrorResponse{
			Error: err1.Error(),
		}
		WriteJsonResponseWithCode(w, response, http.StatusBadRequest)
		return
	}
	res, err2 := app.CoreApp.HandleRelay(relay)
	if err2 != nil {
		response := types.RpcRelayErrorResponse{
			Error: err2.Error(),
		}
		WriteJsonResponseWithCode(w, response, http.StatusBadRequest)
		return
	}
	if logger.DebugEnabled() {
		logger.Logger().Debug().Str("response", fmt.Sprintf("%#v", res)).Msg("relay response")
	}
	WriteResponse(w, res)
}

func extractRelay(w http.ResponseWriter, r *http.Request, p httprouter.Params) (coretypes.Relay, error) {
	relay := coretypes.Relay{}
	body, err1 := ExtractBody(w, r, p)
	if err1 != nil {
		return relay, err1
	}
	err2 := ExtractModelFromBody(body, relay)
	if err2 != nil {
		return relay, err2
	}
	// If model could not be extracted from body -> set default
	if relay.Chain == "" {
		relay = coretypes.Relay{
			Payload: coretypes.Payload{
				Data:    string(body),
				Method:  http.MethodPost,
				Path:    "",
				Headers: r.Header,
			},
		}
	}
	if relay.Host == "" {
		relay.Host = r.Host
	}
	return relay, nil
}