package types

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/shared/pkg/util/jsonutil"
)

func TestJsonMarshal(t *testing.T) {
	t.Parallel()
	node := &Node{
		Id: "1234",
		LastUpdate: 123, // should be serialized as string
	}
	expected := `{"id":"1234","lastUpdate":"123","delayTime":"0","lastBlock":{},"httpInfo":{}}`
	actual := jsonutil.ToString(node)
	require.Equal(t, expected, actual)
}

func TestJsonUnmarshal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		json string
		expected *Node
	}{
		{
			name: "Test unmarshal string",
			json: `{"id":"1234", "lastUpdate":"123"}`,
			expected: &Node{
				Id: "1234",
				LastUpdate: 123,
			},
		},
		{
			name: "Test unmarshal number",
			json: `{"id":"1234", "lastUpdate":123}`,
			expected: &Node{
				Id: "1234",
				LastUpdate: 123,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var actual Node
			err := jsonutil.FromByte([]byte(tc.json), &actual)
			require.NoError(t, err)
			require.Equal(t, tc.expected.Id, actual.Id)
			require.Equal(t, tc.expected.LastUpdate, actual.LastUpdate)
		})
	}
}

func TestJsonUnmarshalError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		json string
		err error
	}{
		{
			name: "Test unmarshal empty",
			json: ``,
			err: errors.New("unexpected end of JSON input"),
		},
		{
			name: "Test unmarshal invalid 1",
			json: `{"lastUpdate":"abcd"}`,
			err: errors.New("invalid character 'a' looking for beginning of value"),
		},
		{
			name: "Test unmarshal invalid 2",
			json: `{"lastUpdate":abcd}`,
			err: errors.New("invalid character 'a' looking for beginning of value"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var actual Node
			err := jsonutil.FromByte([]byte(tc.json), &actual)
			require.Error(t, err)
			require.Equal(t, tc.err.Error(), err.Error())
		})
	}
}