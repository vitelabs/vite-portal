package main

import (
	"fmt"
	"time"

	"github.com/vitelabs/vite-portal/shared/pkg/ws"
)

var timeout = 1000 * time.Millisecond

func main() {
	mock := ws.NewMockWsRpc(0)
	fmt.Println(mock.Url)
	mock.Serve(timeout)
}
