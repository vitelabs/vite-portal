package rpc

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/vitelabs/vite-portal/internal/app"
	nodetypes "github.com/vitelabs/vite-portal/internal/node/types"
)

func GetChains(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	res := app.CoreApp.GetChains()
	WriteResponse(w, res)
}

func GetNodes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var params = nodetypes.GetNodesParams{}
	if err := ExtractModel(w, r, p, &params); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := app.CoreApp.GetNodes(params)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
	WriteResponse(w, res)
}

func PutNode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func DeleteNode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	
}