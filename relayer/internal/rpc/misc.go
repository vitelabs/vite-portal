package rpc

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/vitelabs/vite-portal/internal/app"
	"github.com/vitelabs/vite-portal/internal/logger"
)

func Name(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//go doTestTask(r.Context())
	WriteResponse(w, app.AppName)
}

func Version(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	//test()
	WriteResponse(w, app.AppVersion)
}

func test() {
	ctx, cancelFn := context.WithCancel(context.Background())
	go cancelTestTask(cancelFn)
	go doTestTask(ctx)
}

func cancelTestTask(cancelFn context.CancelFunc) {
	time.Sleep(time.Second * 5)
	cancelFn()
}

func doTestTask(ctx context.Context) {
	defer logger.Logger().Info().Msg("context cancelled")
	ticker := time.NewTicker(1 * time.Second)
	for {
			select {
			case <-ticker.C:
				logger.Logger().Info().Msg(fmt.Sprintf("test: %d", time.Now().UnixMilli()))
			case <-ctx.Done():
					return
			}
	}	
}