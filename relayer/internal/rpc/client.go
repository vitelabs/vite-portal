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
	if app.GlobalConfig.Debug {
		logger.Logger().Debug().Str("request", fmt.Sprintf("%#v", r)).Msg("relay request")
	}
	body, err1 := ExtractBody(w, r, p)
	if err1 != nil {
		response := types.RpcRelayErrorResponse{
			Error: err1,
		}
		WriteResponseWithCode(w, response, http.StatusBadRequest)
		return
	}
	var relay = coretypes.Relay{
		Payload: coretypes.Payload{
			Data:    string(body),
			Method:  r.Method,
			Path:    "",
			Headers: r.Header,
		},
	}
	res, err2 := app.CoreApp.HandleRelay(relay)
	if err2 != nil {
		response := types.RpcRelayErrorResponse{
			Error:    err2,
		}
		WriteResponseWithCode(w, response, http.StatusBadRequest)
		return
	}
	WriteResponse(w, res)
}