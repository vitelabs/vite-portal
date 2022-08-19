package rpc

import (
	"crypto/rand"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

func TestWsConnection(t *testing.T) {
	ws := startWsRpc(types.DefaultRpcWsPort)
	writeMsg := []byte("hello, world!\n")
	if err := ws.WriteMessage(websocket.BinaryMessage, writeMsg); err != nil {
		log.Fatal(err)
	}
	_, n, err := ws.ReadMessage()
	require.NoError(t, err)
	fmt.Printf("Received: %s.\n", n)
}

func TestWsMaxPayloadBytes(t *testing.T) {
	ws := startWsRpc(types.DefaultRpcWsPort + 1)
	writeMsg := make([]byte, sharedtypes.MaxPayloadSize + 1)
	rand.Read(writeMsg)
	err := ws.WriteMessage(websocket.BinaryMessage, writeMsg)
	require.NoError(t, err)

	_, _, err = ws.ReadMessage()
	require.Error(t, err)
	require.Equal(t, "websocket: close 1009 (message too big)", err.Error())

	err = ws.WriteMessage(websocket.TextMessage, []byte("test"))
	require.Error(t, err)
	require.Equal(t, "websocket: close sent", err.Error())
}

func startWsRpc(port int32) *websocket.Conn {
	go StartWsRpc(port, 0)
	time.Sleep(time.Duration(time.Millisecond * 100))
	url := fmt.Sprintf("ws://localhost:%d/ws/v1", port)
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}
	return ws
}