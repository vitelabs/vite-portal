package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/idutil"
)

func TestConnect(t *testing.T) {
	t.Skip()
	c := NewClient("ws://localhost:57331/", sharedtypes.DefaultJwtSecret, time.Duration(types.DefaultRpcTimeout) * time.Millisecond, 0)
	err := c.Connect(idutil.NewGuid())
	require.NoError(t, err)
}

func TestCreateToken(t *testing.T) {
	c := NewClient("ws://localhost:57331/", sharedtypes.DefaultJwtSecret, time.Duration(types.DefaultRpcTimeout) * time.Millisecond, 0)
	token := c.CreateToken("test1234", 0) // never expires
	fmt.Println(token)
	require.NotEmpty(t, token)
}