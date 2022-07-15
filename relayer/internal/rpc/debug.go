package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/julienschmidt/httprouter"
)

func debugTest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	i := 1
	ticker := time.NewTicker(1 * time.Second)
	for {
			select {
			case <-ticker.C:
				fmt.Println(i, "debugTest")
				i++
			case <-r.Context().Done():
					fmt.Println(i, "debugTest done")
					return
			}
	}
}

func debugMemStats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	b, err := json.Marshal(m)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
	}
	WriteResponse(w, string(b), r.URL.Path, r.Host)
}