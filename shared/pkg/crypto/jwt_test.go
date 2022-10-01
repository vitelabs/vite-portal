package crypto

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/shared/pkg/types"
)

func TestJwt_GetClaims(t *testing.T) {
	h := NewJWTHandler([]byte("secret"), time.Minute)
	token := h.IssueDefaultToken("test_subject", "test_issuer", 0)
	header := map[string][]string{types.HTTPHeaderAuthorization: {fmt.Sprintf("Bearer %s", token)}}
	claims, err := h.GetClaims(header)
	require.NoError(t, err)
	require.NotNil(t, claims)
	require.Equal(t, "test_subject", claims.Subject)
	require.Equal(t, "test_issuer", claims.Issuer)
	require.Empty(t, claims.ExpiresAt)
}

func TestJwt_GetClaims_ExpiresAt(t *testing.T) {
	h := NewJWTHandler([]byte("secret"), time.Minute)
	now := time.Now().Unix()
	token := h.IssueDefaultToken("test_subject", "test_issuer", 1)
	header := map[string][]string{types.HTTPHeaderAuthorization: {fmt.Sprintf("Bearer %s", token)}}
	claims, err := h.GetClaims(header)
	require.NoError(t, err)
	require.NotNil(t, claims)
	require.Equal(t, now + 1, claims.ExpiresAt.Unix())
}

func TestJwt_GetClaims_Ivalid(t *testing.T) {
	h := NewJWTHandler([]byte("secret"), time.Minute)
	header := map[string][]string{types.HTTPHeaderAuthorization: {"Bearer token1234"}}
	claims, err := h.GetClaims(header)
	require.Error(t, err)
	require.NotNil(t, claims)
	require.Equal(t, "", claims.Subject)
}

func TestJwt_GetClaims_Empty(t *testing.T) {
	h := NewJWTHandler([]byte("secret"), time.Minute)
	claims, err := h.GetClaims(nil)
	require.Error(t, err)
	require.NotNil(t, claims)
	require.Equal(t, "", claims.Subject)
}
