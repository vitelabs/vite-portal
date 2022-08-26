package app

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	coretypes "github.com/vitelabs/vite-portal/relayer/internal/core/types"
	nodetypes "github.com/vitelabs/vite-portal/relayer/internal/node/types"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/httputil"
)

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

func (a *RelayerApp) GetChains(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	res := a.nodeService.GetChains()
	httputil.WriteJsonResponse(w, res)
}

func (a *RelayerApp) GetNodes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var params = nodetypes.GetNodesParams{}
	if err := httputil.ExtractQuery(w, r, &params); err != nil {
		httputil.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	params.Offset, params.Limit = checkPagination(params.Offset, params.Limit)
	res, err := a.nodeService.GetNodes(params.Chain, params.Offset, params.Limit)
	if err != nil {
		httputil.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.WriteJsonResponse(w, res)
}

func (a *RelayerApp) GetNode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if id == "" {
		httputil.WriteErrorResponse(w, http.StatusBadRequest, "invalid identifier")
		return
	}
	res, found := a.nodeService.GetNode(id)
	if !found {
		httputil.WriteErrorResponse(w, http.StatusNotFound, "node does not exist")
		return
	}
	httputil.WriteJsonResponse(w, res)
}

// PutNode enables the orchestrator to add or update a node
func (a *RelayerApp) PutNode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: auth
	var node = nodetypes.Node{}
	if err := httputil.ExtractModel(w, r, &node, sharedtypes.MaxPayloadSize); err != nil {
		httputil.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := a.nodeService.PutNode(node); err != nil {
		httputil.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	httputil.WriteJsonResponse(w, nil)
}

// DeleteNode enables the orchestrator to delete a node
func (a *RelayerApp) DeleteNode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: auth
	id := p.ByName("id")
	if id == "" {
		httputil.WriteErrorResponse(w, http.StatusBadRequest, "invalid identifier")
		return
	}
	if err := a.nodeService.DeleteNode(id); err != nil {
		httputil.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	httputil.WriteJsonResponse(w, nil)
}

func checkPagination(offset, limit int) (int, int) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 1000 {
		limit = 1000
	}
	return offset, limit
}