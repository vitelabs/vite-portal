package rpc

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/httputil"
	"github.com/vitelabs/vite-portal/shared/pkg/version"
)

func Name(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	httputil.WriteResponse(w, types.AppName, httputil.ContentTypeTextPlain)
}

func Version(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	httputil.WriteResponse(w, version.PROJECT_BUILD_VERSION, httputil.ContentTypeTextPlain)
}
