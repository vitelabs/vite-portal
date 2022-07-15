package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCodeTypeMessages(t *testing.T) {
	expected := codeTypeLimit
	actual := len(CodeTypeMessages)
	require.Equal(t, actual, expected)
}

func TestUnknownCodeTypeMessage(t * testing.T) {
	expected := "unknown code 1000"
	actual := GetCodeMessage(1000)
	require.Equal(t, actual, expected)
}

func TestCodeTypeName(t *testing.T) {
	expected1 := "everything ok"
	actual1 := GetCodeMessage(CodeOK)
	require.Equal(t, actual1, expected1)

	expected2 := "internal error"
	actual2 := GetCodeMessage(CodeInternal)
	require.Equal(t, actual2, expected2)
}