package rpc

import (
	"crypto/rand"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/internal/types"
	"golang.org/x/net/websocket"
)

func TestWsConnection(t *testing.T) {
	ws := startWsRpc(types.DefaultRpcWsPort)
	writeMsg := []byte("hello, world!\n")
	if _, err := ws.Write(writeMsg); err != nil {
		log.Fatal(err)
	}
	var readMsg = make([]byte, 512)
	n, err := ws.Read(readMsg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s.\n", readMsg[:n])
}

func TestWsMaxPayloadBytes(t *testing.T) {
	ws := startWsRpc(types.DefaultRpcWsPort + 1)
	writeMsg := make([]byte, types.MaxRequestContentLength + 1)
	rand.Read(writeMsg)
	if _, err := ws.Write(writeMsg); err != nil {
		log.Fatal(err)
	}
	var readMsg = make([]byte, 512)
	n, err := ws.Read(readMsg)
	if err != nil {
		log.Fatal(err)
	}
	require.Equal(t, "Code: -32600 Message: websocket: frame payload size exceeds limit", fmt.Sprintf("%s", readMsg[:n]))
}

func startWsRpc(port int32) *websocket.Conn {
	go StartWsRpc(port, 0)
	time.Sleep(time.Duration(time.Millisecond * 100))
	origin := "http://localhost/"
	url := fmt.Sprintf("ws://localhost:%d/ws/v1", port)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	return ws
}