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
	timeout := 100 * time.Millisecond
	mock := StartMockWsRpc(0, timeout)
	require.NotNil(t, mock)
	require.Greater(t, mock.Port, 0)
	require.Equal(t, "/ws/mock", mock.Pattern)
	require.Equal(t, fmt.Sprintf("ws://localhost:%d/ws/mock", mock.Port), mock.Url)
	conn, _, err := websocket.DefaultDialer.Dial(mock.Url, nil)
	if err != nil {
		log.Fatal(err)
	}
	require.NotNil(t, conn)
}