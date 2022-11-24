package relayer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/shared/pkg/crypto"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

func TestIssueDefaultToken(t *testing.T) {
	h := crypto.NewJWTHandler([]byte(sharedtypes.DefaultJwtSecret), 0)
	token := h.IssueDefaultToken(sharedtypes.DefaultJwtSecret, sharedtypes.JWTOrchestratorIssuer, 0) // never expires
	fmt.Println(token)
	require.NotEmpty(t, token)
}