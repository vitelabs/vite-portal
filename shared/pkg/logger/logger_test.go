package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDebugLogger(t *testing.T) {
	tests := []struct {
		name string
		debug bool
	}{
		{
			name: "Test debug=false",
			debug: false,
		},
		{
			name: "Test debug=true",
			debug: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			Init(tc.debug)
			called := false
			logger.Debug().Str("test", testFn(&called))
			// testFn is called in both cases
			require.True(t, called)
		})
	}
}

func testFn(called *bool) string  {
	(*called) = true
	return "test1234"
}