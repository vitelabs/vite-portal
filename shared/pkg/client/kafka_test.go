package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/shared/pkg/types"
)

var defaultKafkaTimeout = 1 * time.Second

func TestKafkaWrite(t *testing.T) {
	t.Skip()
	cfg := types.DefaultKafkaConfig
	c := NewKafkaClient(defaultKafkaTimeout, cfg.Server, cfg.DefaultTopic)
	round := time.Now().UnixMilli() / 1000 / 60
	c.Write(fmt.Sprintf("id: %d", round))
	c.Write(fmt.Sprintf("id: %d (2)", round))
	c.Close()
}

func TestKafkaRead(t *testing.T) {
	t.Skip()
	cfg := types.DefaultKafkaConfig
	c := NewKafkaClient(defaultKafkaTimeout, cfg.Server, cfg.DefaultTopic)
	fmt.Println("Round 1:")
	messages, err := c.Read(0, 2, defaultKafkaTimeout)
	require.NoError(t, err)
	for _, m := range messages {
		fmt.Println(m)
	}
	fmt.Println("Round 2:")
	messages, err = c.Read(2, 2, defaultKafkaTimeout)
	require.NoError(t, err)
	for _, m := range messages {
		fmt.Println(m)
	}
	fmt.Println("Round 3:")
	messages, err = c.Read(4, 10000, defaultKafkaTimeout)
	require.NoError(t, err)
	for _, m := range messages {
		fmt.Println(m)
	}
	c.Close()
}
