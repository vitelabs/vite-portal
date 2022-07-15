package rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/vitelabs/vite-portal/internal/logger"
	"github.com/vitelabs/vite-portal/internal/types"
)

type route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

func StartHttpRpc(port int32, timeout int64, debug bool, profile bool) {
	routes := []route{
		{Name: "Default", Method: "GET", Path: "/", HandlerFunc: Name},
		{Name: "AppName", Method: "GET", Path: "/api", HandlerFunc: Name},
		{Name: "AppVersion", Method: "GET", Path: "/api/v1", HandlerFunc: Version},
		{Name: "Relay", Method: "POST", Path: "/api/v1/client/relay", HandlerFunc: Relay},
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

func WriteResponse(w http.ResponseWriter, jsn, path, ip string) {
	b, err := json.Marshal(jsn)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		fmt.Println(err.Error())
	} else {
		writeDefaultHeader(w)
		_, err := w.Write(b)
		if err != nil {
			logger.Logger().Error().Err(err).Msg("WriteResponse failed")
		}
	}
}

func WriteRaw(w http.ResponseWriter, jsn, path, ip string) {
	writeDefaultHeader(w)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(jsn))
	if err != nil {
		logger.Logger().Error().Err(err).Msg("WriteRaw failed")
	}
}

func WriteJsonResponse(w http.ResponseWriter, jsonData, path, ip string) {
	WriteJsonResponseWithCode(w, jsonData, path, ip, http.StatusOK)
}

func WriteJsonResponseWithCode(w http.ResponseWriter, jsonData, path, ip string, code int) {
	writeDefaultHeader(w)
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &raw); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		logger.Logger().Error().Err(err).Msg("WriteJsonResponseWithCode Unmarshal failed")
		return
	}
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(raw)
	if err != nil {
		logger.Logger().Error().Err(err).Msg("WriteJsonResponseWithCode Encode failed")
		return
	}
}

func WriteErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	writeDefaultHeader(w)
	w.WriteHeader(errorCode)
	err := json.NewEncoder(w).Encode(&types.RpcError{
		Code:    errorCode,
		Message: errorMsg,
	})
	if err != nil {
		logger.Logger().Error().Err(err).Msg("WriteErrorResponse failed")
	}
}

func writeDefaultHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func ExtractModel(_ http.ResponseWriter, r *http.Request, _ httprouter.Params, model interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // 1048576 bytes = 1 megabyte
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return nil
	}
	if err := r.Body.Close(); err != nil {
		return err
	}
	if err := json.Unmarshal(body, model); err != nil {
		return err
	}
	return nil
}