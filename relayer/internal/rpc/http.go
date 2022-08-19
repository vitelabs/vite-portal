package rpc

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/jsonutil"
)

const (
	ContentTypeTextPlain = "text/plain; charset=UTF-8"
	ContentTypeJson      = "application/json; charset=UTF-8"
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

func WriteResponse(w http.ResponseWriter, data, contentType string) {
	WriteResponseWithCode(w, data, contentType, http.StatusOK)
}

func WriteResponseWithCode(w http.ResponseWriter, data, contentType string, code int) {
	writeHeader(w, contentType)
	w.WriteHeader(code)
	_, err2 := w.Write([]byte(data))
	if err2 != nil {
		logger.Logger().Error().Err(err2).Msg("WriteResponseWithCode failed")
	}
}

func WriteJsonResponse(w http.ResponseWriter, data any) {
	WriteJsonResponseWithCode(w, data, http.StatusOK)
}

func WriteJsonResponseWithCode(w http.ResponseWriter, data any, code int) {
	b, err1 := jsonutil.ToByte(data)
	if err1 != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err1.Error())
		logger.Logger().Error().Err(err1).Msg("WriteJsonResponseWithCode failed")
		return
	}
	writeHeader(w, ContentTypeJson)
	w.WriteHeader(code)
	_, err2 := w.Write(b)
	if err2 != nil {
		logger.Logger().Error().Err(err2).Msg("WriteJsonResponseWithCode failed")
	}
}

func WriteErrorResponse(w http.ResponseWriter, code int, msg string) {
	err := &types.RpcError{
		Code:    code,
		Message: msg,
	}
	WriteJsonResponseWithCode(w, err, code)
}

func writeHeader(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType)
}
