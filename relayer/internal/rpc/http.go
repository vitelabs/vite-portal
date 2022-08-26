package rpc

import (
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
)

type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

func StartHttpRpc(routes []Route, port int32, timeout time.Duration, debug, profile bool) {
	routes = append(routes, []Route{
		{Name: "Default", Method: "GET", Path: "/", HandlerFunc: Name},
		{Name: "AppName", Method: "GET", Path: "/api", HandlerFunc: Name},
		{Name: "AppVersion", Method: "GET", Path: "/api/v1", HandlerFunc: Version},
	}...)

	if debug {
		routes = append(routes, Route{Name: "DebugTest", Method: "GET", Path: "/debug/test", HandlerFunc: debugTest})
	}

	if profile {
		routes = append(routes, Route{Name: "ProfileMemStats", Method: "GET", Path: "/profile/memstats", HandlerFunc: debugMemStats})
	}

	srv := &http.Server{
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		Addr:              ":" + strconv.Itoa(int(port)),
		Handler:           http.TimeoutHandler(router(routes), timeout, "Server Timeout Handling Request"),
	}

	l, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		logger.Logger().Error().Err(err).Msg("HTTP RPC error")
		os.Exit(1)
	}

	go func() {
		err := srv.Serve(l)
		if err != nil {
			logger.Logger().Error().Err(err).Msg("HTTP RPC closed")
			os.Exit(1)
		}
	}()
}

func router(routes []Route) *httprouter.Router {
	router := httprouter.New()
	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.HandlerFunc)
	}
	return router
}
