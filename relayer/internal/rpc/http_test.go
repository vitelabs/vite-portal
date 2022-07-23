package rpc

import (
	"errors"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
	nodetypes "github.com/vitelabs/vite-portal/internal/node/types"
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
			if err != nil || tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, err.Error(), tc.expectedError.Error())
			}
			require.Equal(t, tc.expected, tc.model)
		})
	}
}

func TestExtractParams(t *testing.T) {
	tests := []struct {
		name          string
		params        httprouter.Params
		model         nodetypes.GetNodesParams
		expected      nodetypes.GetNodesParams
		expectedError error
	}{
		{
			name:          "Test nil params",
			params:        nil,
			model:         nodetypes.GetNodesParams{},
			expected:      nodetypes.GetNodesParams{},
			expectedError: nil,
		},
		{
			name: "Test invalid params",
			params: httprouter.Params{
				httprouter.Param{Key: "asdf", Value: "1234"},
			},
			model:         nodetypes.GetNodesParams{},
			expected:      nodetypes.GetNodesParams{},
			expectedError: nil,
		},
		{
			name: "Test chain only",
			params: httprouter.Params{
				httprouter.Param{Key: "chain", Value: "chain1"},
			},
			model: nodetypes.GetNodesParams{},
			expected: nodetypes.GetNodesParams{
				Chain: "chain1",
			},
			expectedError: nil,
		},
		{
			name: "Test 2 chainz",
			params: httprouter.Params{
				httprouter.Param{Key: "chain", Value: "chain1"},
				httprouter.Param{Key: "chain", Value: "chain2"},
			},
			model: nodetypes.GetNodesParams{},
			expected: nodetypes.GetNodesParams{
				Chain: "chain2",
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ExtractParams(nil, nil, tc.params, &tc.model)
			if err != nil || tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, err.Error(), tc.expectedError.Error())
			}
			require.Equal(t, tc.expected, tc.model)
		})
	}
}
