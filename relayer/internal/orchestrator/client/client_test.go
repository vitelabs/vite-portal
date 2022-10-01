package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/idutil"
)

func TestConnect(t *testing.T) {
	t.Skip()
	c := NewClient("ws://localhost:57331/", types.DefaultJwtSecret, time.Duration(types.DefaultRpcTimeout) * time.Millisecond)
	err := c.Connect(idutil.NewGuid())
	require.NoError(t, err)
}

func TestCreateToken(t *testing.T) {
	c := NewClient("ws://localhost:57331/", types.DefaultJwtSecret, time.Duration(types.DefaultRpcTimeout) * time.Millisecond)
	token := c.CreateToken("test1234", 28800) // expires in 8 hours
	fmt.Println(token)
	require.NotEmpty(t, token)
}