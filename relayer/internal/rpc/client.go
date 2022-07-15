package rpc

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/vitelabs/vite-portal/internal/app"
	coretypes "github.com/vitelabs/vite-portal/internal/core/types"
	"github.com/vitelabs/vite-portal/internal/types"
)

func Relay(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if cors(&w, r) {
		return
	}
	var relay = coretypes.Relay{}
	if err1 := ExtractModel(w, r, ps, &relay); err1 != nil {
		response := types.RpcRelayErrorResponse{
			Error: err1,
		}
		j, _ := json.Marshal(response)
		WriteJsonResponseWithCode(w, string(j), r.URL.Path, r.Host, http.StatusBadRequest)
		return
	}
	res, err2 := app.CoreApp.HandleRelay(relay)
	if err2 != nil {
		response := types.RpcRelayErrorResponse{
			Error:    err2,
		}
		j, _ := json.Marshal(response)
		WriteJsonResponseWithCode(w, string(j), r.URL.Path, r.Host, http.StatusBadRequest)
		return
	}
	j, err3 := json.Marshal(res)
	if err3 != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err3.Error())
		return
	}
	WriteJsonResponse(w, string(j), r.URL.Path, r.Host)
}