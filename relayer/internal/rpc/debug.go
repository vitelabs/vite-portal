package rpc

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/vitelabs/vite-portal/shared/pkg/util/httputil"
)

func debugTest(_ http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

func debugMemStats(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	httputil.WriteJsonResponse(w, m)
}