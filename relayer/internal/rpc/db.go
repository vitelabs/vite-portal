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
	if err := ExtractQuery(w, r, p, &params); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := app.CoreApp.GetNodes(params)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
	WriteResponse(w, res)
}

func GetNode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if id == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "invalid identifier")
		return
	}
	res, found := app.CoreApp.GetNode(id)
	if !found {
		WriteErrorResponse(w, http.StatusNotFound, "node does not exist")
		return
	}
	WriteResponse(w, res)
}

// PutNode enables the orchestrator to add or update a node
func PutNode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: auth
	var node = nodetypes.Node{}
	if err := ExtractModel(w, r, p, &node); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := app.CoreApp.PutNode(node); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	WriteResponse(w, nil)
}

// DeleteNode enables the orchestrator to delete a node
func DeleteNode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: auth
}
