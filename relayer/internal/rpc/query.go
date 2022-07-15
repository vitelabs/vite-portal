package rpc

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/vitelabs/vite-portal/internal/app"
)

func Name(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	WriteResponse(w, app.AppName, r.URL.Path, r.Host)
}

func Version(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	WriteResponse(w, app.AppVersion, r.URL.Path, r.Host)
}