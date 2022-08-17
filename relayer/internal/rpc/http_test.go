package rpc

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	nodetypes "github.com/vitelabs/vite-portal/relayer/internal/node/types"
)

func TestExtractModelFromBody(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		model         nodetypes.GetNodesParams
		expected      nodetypes.GetNodesParams
		expectedError error
	}{
		{
			name:          "Test empty body",
			body:          "",
			model:         nodetypes.GetNodesParams{},
			expected:      nodetypes.GetNodesParams{},
			expectedError: nil,
		},
		{
			name:          "Test invalid body",
			body:          "1234",
			model:         nodetypes.GetNodesParams{},
			expected:      nodetypes.GetNodesParams{},
			expectedError: errors.New("json: cannot unmarshal number into Go value of type types.GetNodesParams"),
		},
		{
			name:  "Test chain only",
			body:  "{ \"chain\": \"chain1\"}",
			model: nodetypes.GetNodesParams{},
			expected: nodetypes.GetNodesParams{
				Chain: "chain1",
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ExtractModelFromBody([]byte(tc.body), &tc.model)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, err.Error(), tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expected, tc.model)
		})
	}
}

func TestExtractParams(t *testing.T) {
	tests := []struct {
		name          string
		request       http.Request
		model         nodetypes.GetNodesParams
		expected      nodetypes.GetNodesParams
		expectedError error
	}{
		{
			name: "Test nil query",
			request: http.Request{
				URL: parseUrl(t, ""),
			},
			model:         nodetypes.GetNodesParams{},
			expected:      nodetypes.GetNodesParams{},
			expectedError: nil,
		},
		{
			name: "Test invalid query",
			request: http.Request{
				URL: parseUrl(t, "http://localhost/v1?test=1234"),
			},
			model:         nodetypes.GetNodesParams{},
			expected:      nodetypes.GetNodesParams{},
			expectedError: nil,
		},
		{
			name: "Test chain only",
			request: http.Request{
				URL: parseUrl(t, "http://localhost/v1?chain=chain1"),
			},
			model: nodetypes.GetNodesParams{},
			expected: nodetypes.GetNodesParams{
				Chain: "chain1",
			},
			expectedError: nil,
		},
		{
			name: "Test 2 chainz",
			request: http.Request{
				URL: parseUrl(t, "http://localhost/v1?chain=chain1&chain=chain2"),
			},
			model: nodetypes.GetNodesParams{},
			expected: nodetypes.GetNodesParams{
				Chain: "chain1",
			},
			expectedError: nil,
		},
		{
			name: "Test full",
			request: http.Request{
				URL: parseUrl(t, "http://localhost/v1?chain=chain1&offset=10&limit=5"),
			},
			model: nodetypes.GetNodesParams{},
			expected: nodetypes.GetNodesParams{
				Chain:  "chain1",
				Offset: 10,
				Limit:  5,
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ExtractQuery(nil, &tc.request, nil, &tc.model)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expected, tc.model)
		})
	}
}

func parseUrl(t *testing.T, raw string) *url.URL {
	p, err := url.Parse(raw)
	require.NoError(t, err)
	return p
}
