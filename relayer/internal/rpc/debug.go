package rpc

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/julienschmidt/httprouter"
)

func debugMemStats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	b, err := json.Marshal(m)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
	}
	WriteResponse(w, string(b), r.URL.Path, r.Host)
}