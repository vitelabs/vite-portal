package rpc

import (
	"errors"
	"log"
	"net/http"

	"github.com/vitelabs/vite-portal/orchestrator/internal/app"
	"github.com/vitelabs/vite-portal/shared/pkg/util/httputil"
)

func handleNode(w http.ResponseWriter, r *http.Request) {
	chain, err := getChain(r)
	if err != nil {
		log.Print(err)
		httputil.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Println(chain)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func getChain(r *http.Request) (string, error) {
	chain := r.URL.Query().Get("chain")
	if chain == "" {
		return app.CoreApp.Config.DefaultChain, nil
	}
	for _, v := range app.CoreApp.Config.SupportedChains {
		if chain == v {
			return chain, nil
		}
	}
	return "", errors.New("chain not supported")
}