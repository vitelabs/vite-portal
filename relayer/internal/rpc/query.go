package rpc

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/vitelabs/vite-portal/internal/app"
	nodetypes "github.com/vitelabs/vite-portal/internal/node/types"
)

func Name(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	WriteResponse(w, app.AppName)
}

func Version(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	WriteResponse(w, app.AppVersion)
}

func Chains(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	res := app.CoreApp.QueryChains()
	WriteResponse(w, res)
}

func Nodes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var params = nodetypes.QueryNodesParams{}
	if err := ExtractModel(w, r, p, &params); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := app.CoreApp.QueryNodes(params)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
	WriteResponse(w, res)
}