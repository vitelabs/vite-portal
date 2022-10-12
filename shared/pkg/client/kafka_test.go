package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/vitelabs/vite-portal/shared/pkg/types"
)

var defaultKafkaTimeout = 1 * time.Second

func TestKafkaWrite(t *testing.T) {
	t.Skip()
	c := NewKafkaClient(defaultKafkaTimeout, *types.NewDefaultKafkaConfig())
	round := time.Now().UnixMilli() / 1000 / 60
	c.Write(fmt.Sprintf("id: %d", round))
	c.Write(fmt.Sprintf("id: %d (2)", round))
	c.Close()
}

func TestKafkaRead(t *testing.T) {
	t.Skip()
	c := NewKafkaClient(defaultKafkaTimeout, *types.NewDefaultKafkaConfig())
	fmt.Println("Round 1:")
	messages := c.Read()
	for _, m := range messages {
		fmt.Println(m)
	}
	fmt.Println("Round 2:")
	messages = c.Read()
	for _, m := range messages {
		fmt.Println(m)
	}
	fmt.Println("Round 3:")
	messages = c.Read()
	for _, m := range messages {
		fmt.Println(m)
	}
	c.Close()
}
