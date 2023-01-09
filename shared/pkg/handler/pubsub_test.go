package handler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/commonutil"
)

func TestPubSubTopics(t *testing.T) {
	timeout := 2 * time.Second
	pb := NewPubsubTopics()
	t1 := "topic1"
	c1 := pb.Subscribe(t1)
	c2 := pb.Subscribe(t1)
	v1 := ""
	v2 := ""
	go func() {
		for {
			select {
			case val := <-c1:
				logger.Logger().Info().Msg(val)
				v1 = val
			}
		}
	}()
	go func() {
		for {
			select {
			case val := <-c2:
				logger.Logger().Info().Msg(val)
				v2 = val
			}
		}
	}()
	expected1 := "val1"
	pb.Publish(t1, expected1)
	require.Equal(t, "", v1)
	require.Equal(t, "", v2)
	commonutil.WaitFor(timeout, func() bool {
		return v2 == expected1
	})
	require.Equal(t, expected1, v1)
	require.Equal(t, expected1, v2)
	expected2 := "val2"
	pb.Publish(t1, expected2)
	require.Equal(t, expected1, v1)
	require.Equal(t, expected1, v2)
	commonutil.WaitFor(timeout, func() bool {
		return v2 == expected2
	})
	require.Equal(t, expected2, v1)
	require.Equal(t, expected2, v2)
}