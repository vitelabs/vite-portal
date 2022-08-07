package rpc

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/vitelabs/vite-portal/internal/app"
	"github.com/vitelabs/vite-portal/version"
)

func Name(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	WriteResponse(w, app.AppName)
}

func Version(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	WriteResponse(w, version.RELAYER_BUILD_VERSION)
}
