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
	c.Start()
	round := time.Now().UnixMilli() / 1000 / 60
	c.Write(fmt.Sprintf("id: %d", round))
	c.Stop()
}

func TestKafkaRead(t *testing.T) {
	t.Skip()
	c := NewKafkaClient(defaultKafkaTimeout, *types.NewDefaultKafkaConfig())
	c.Start()
	messages := c.Read()
	for _, m := range messages {
		fmt.Println(m)
	}
	c.Stop()
}
