package types

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	roottypes "github.com/vitelabs/vite-portal/relayer/internal/types"
)

func TestHttpExecutionError(t *testing.T) {
	t.Parallel()
	err := NewError(DefaultCodeNamespace, CodeHttpExecutionError, errors.New("test1234"))

	require.NotNil(t, err)
	require.Equal(t, int(err.Code()), CodeHttpExecutionError)
	require.Equal(t, err.CodeNamespace(), DefaultCodeNamespace)
	require.Equal(t, reflect.TypeOf(err.Data()), reflect.TypeOf(roottypes.FmtError{}))
	expectedError := `ERROR:
Namespace: core
Code: 1
Message: %s
`
	require.Equal(t, fmt.Sprintf(expectedError, err.Error()), err.ErrorFormatted())

	data, ok := err.Data().(roottypes.FmtError)
	require.True(t, ok)
	require.Equal(t, "error executing the http request: test1234", data.Error())
}
