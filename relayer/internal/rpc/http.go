package rpc

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
)

type route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

func StartHttpRpc(port int32, timeout int64, debug, profile bool) {
	routes := []route{
		{Name: "Default", Method: "GET", Path: "/", HandlerFunc: Name},
		{Name: "AppName", Method: "GET", Path: "/api", HandlerFunc: Name},
		{Name: "AppVersion", Method: "GET", Path: "/api/v1", HandlerFunc: Version},
		{Name: "Relay", Method: "POST", Path: "/api/v1/client/relay", HandlerFunc: Relay},
		{Name: "GetChains", Method: "GET", Path: "/api/v1/db/chains", HandlerFunc: GetChains},
		{Name: "GetNodes", Method: "GET", Path: "/api/v1/db/nodes", HandlerFunc: GetNodes},
		{Name: "GetNode", Method: "GET", Path: "/api/v1/db/nodes/:id", HandlerFunc: GetNode},
		{Name: "PutNode", Method: "PUT", Path: "/api/v1/db/nodes", HandlerFunc: PutNode},
		{Name: "DeleteNode", Method: "DELETE", Path: "/api/v1/db/nodes/:id", HandlerFunc: DeleteNode},
	}

	if debug {
		routes = append(routes, route{Name: "DebugTest", Method: "GET", Path: "/debug/test", HandlerFunc: debugTest})
	}

	if profile {
		routes = append(routes, route{Name: "ProfileMemStats", Method: "GET", Path: "/profile/memstats", HandlerFunc: debugMemStats})
	}

	srv := &http.Server{
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		Addr:              ":" + strconv.Itoa(int(port)),
		Handler:           http.TimeoutHandler(router(routes), time.Duration(timeout)*time.Millisecond, "Server Timeout Handling Request"),
	}

	err := srv.ListenAndServe()
	if err != nil {
		logger.Logger().Error().Err(err).Msg("HTTP RPC closed")
		os.Exit(1)
	}
}

func router(routes []route) *httprouter.Router {
	router := httprouter.New()
	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.HandlerFunc)
	}
	return router
}

func cors(w *http.ResponseWriter, r *http.Request) (isOptions bool) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST")
	// TODO: set Access-Control-Allow-Headers
	return ((*r).Method == "OPTIONS")
}
