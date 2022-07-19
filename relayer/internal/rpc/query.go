package rpc

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/vitelabs/vite-portal/internal/app"
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