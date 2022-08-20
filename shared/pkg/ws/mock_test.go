package ws

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestStartMockWsRpc(t *testing.T) {
	mock := StartMockWsRpc(1000)
	require.Greater(t, mock.Port, 0)
	require.Equal(t, "/ws/mock", mock.Pattern)
	require.Equal(t, fmt.Sprintf("ws://localhost:%d/ws/mock", mock.Port), mock.Url)

	time.Sleep(time.Duration(time.Millisecond * 100))
	conn, _, err := websocket.DefaultDialer.Dial(mock.Url, nil)
	if err != nil {
		log.Fatal(err)
	}
	require.NotNil(t, conn)
}