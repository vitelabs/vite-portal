package jsonutil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testEntry struct {
	Id    string `json:"id"`
	Chain string `json:"chain"`
}

func TestFromByte(t *testing.T) {
	json := `{"id": "1234", "chain": "chain1"}`
	expected := testEntry{
		Id:    "1234",
		Chain: "chain1",
	}
	entry := testEntry{}
	err := FromByte([]byte(json), &entry)
	assert.NoError(t, err)
	assert.Equal(t, "1234", entry.Id)
	assert.Equal(t, "chain1", entry.Chain)
	assert.Equal(t, expected, entry)
}

func TestFromByteMultiple(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		model         testEntry
		expected      testEntry
		expectedError error
	}{
		{
			name:          "Test empty entry",
			body:          "",
			model:         testEntry{},
			expected:      testEntry{},
			expectedError: errors.New("unexpected end of JSON input"),
		},
		{
			name:  "Test chain only",
			body:  "{ \"chain\": \"chain1\"}",
			model: testEntry{},
			expected: testEntry{
				Chain: "chain1",
			},
			expectedError: nil,
		},
		{
			name:  "Test entry",
			body:  "{ \"id\": \"1234\", \"chain\": \"chain1\"}",
			model: testEntry{},
			expected: testEntry{
				Id:    "1234",
				Chain: "chain1",
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := FromByte([]byte(tc.body), &tc.model)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.expected, tc.model)
		})
	}
}
