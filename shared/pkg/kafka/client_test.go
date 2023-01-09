package kafka

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/shared/pkg/types"
)

var defaultTimeout = 2 * time.Second

func TestRead(t *testing.T) {
	t.Skip()
	cfg := types.DefaultKafkaConfig
	c := NewClient(defaultTimeout, cfg.Server, cfg.DefaultTopic)
	fmt.Println("Round 1:")
	messages, err := c.Read(0, 2, defaultTimeout)
	require.NoError(t, err)
	for _, m := range messages {
		fmt.Println(m)
	}
	fmt.Println("Round 2:")
	messages, err = c.Read(2, 2, defaultTimeout)
	require.NoError(t, err)
	for _, m := range messages {
		fmt.Println(m)
	}
	fmt.Println("Round 3:")
	messages, err = c.Read(4, 10000, defaultTimeout)
	require.NoError(t, err)
	for _, m := range messages {
		fmt.Println(m)
	}
	c.Close()
}

func TestWrite(t *testing.T) {
	t.Skip()
	cfg := types.DefaultKafkaConfig
	c := NewClient(defaultTimeout, cfg.Server, cfg.DefaultTopic)
	round := time.Now().UnixMilli() / 1000 / 60
	c.Write(fmt.Sprintf("id: %d", round), "a")
	c.Write(fmt.Sprintf("id: %d (2)", round), "b")
	c.Close()
}
