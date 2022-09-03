package app

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	coretypes "github.com/vitelabs/vite-portal/relayer/internal/core/types"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/httputil"
	"github.com/vitelabs/vite-portal/shared/pkg/version"
)

func (a *RelayerApp) AppInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	res := sharedtypes.RpcAppInfoResponse{
		Id:      a.id,
		Version: version.PROJECT_BUILD_VERSION,
		Name:    types.AppName,
	}
	httputil.WriteJsonResponse(w, res)
}

func (a *RelayerApp) Relay(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
		httputil.WriteJsonResponseWithCode(w, response, http.StatusBadRequest)
		return
	}
	res, err2 := a.HandleRelay(relay)
	if err2 != nil {
		response := types.RpcRelayErrorResponse{
			Error: err2.Error(),
		}
		httputil.WriteJsonResponseWithCode(w, response, http.StatusBadRequest)
		return
	}
	if logger.DebugEnabled() {
		logger.Logger().Debug().Str("response", fmt.Sprintf("%#v", res)).Msg("relay response")
	}
	httputil.WriteResponse(w, res, httputil.ContentTypeJson)
}

func cors(w *http.ResponseWriter, r *http.Request) (isOptions bool) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST")
	// TODO: set Access-Control-Allow-Headers
	return ((*r).Method == "OPTIONS")
}

func extractRelay(w http.ResponseWriter, r *http.Request, p httprouter.Params) (coretypes.Relay, error) {
	relay := coretypes.Relay{}
	body, err1 := httputil.ExtractBody(r, sharedtypes.MaxPayloadSize)
	if err1 != nil {
		return relay, err1
	}
	err2 := httputil.ExtractModelFromBody(body, relay)
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
