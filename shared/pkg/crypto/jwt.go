package crypto

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/types"
)

type JWTHandler struct {
	expiryTimeout time.Duration
	keyFunc       func(token *jwt.Token) (interface{}, error)
}

func NewDefaultJWTHandler(secret []byte) *JWTHandler {
	return NewJWTHandler(secret, types.JWTExpiryTimeout)
}

func NewJWTHandler(secret []byte, expiryTimeout time.Duration) *JWTHandler {
	return &JWTHandler{
		expiryTimeout: expiryTimeout,
		keyFunc: func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		},
	}
}

func (h *JWTHandler) Extract(header http.Header) (string, error) {
	var strToken string
	if auth := header.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") {
		strToken = strings.TrimPrefix(auth, "Bearer ")
	}
	if len(strToken) == 0 {
		return "", errors.New("missing token")
	}
	return strToken, nil
}

func (h *JWTHandler) Validate(token string) (jwt.RegisteredClaims, error) {
	var claims jwt.RegisteredClaims
	// We explicitly set only HS256 allowed, and also disables the
	// claim-check: the RegisteredClaims internally requires 'iat' to
	// be no later than 'now', but we allow for a bit of drift.
	t, err := jwt.ParseWithClaims(token, &claims, h.keyFunc,
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithoutClaimsValidation())

	switch {
	case err != nil:
		return claims, err
	case !t.Valid:
		return claims, errors.New("invalid token")
	case !claims.VerifyExpiresAt(time.Now(), false): // optional
		return claims, errors.New("token is expired")
	case claims.IssuedAt == nil:
		return claims, errors.New("missing issued-at")
	case time.Since(claims.IssuedAt.Time) > h.expiryTimeout:
		return claims, errors.New("stale token")
	case time.Until(claims.IssuedAt.Time) > h.expiryTimeout:
		return claims, errors.New("future token")
	}

	return claims, nil
}

func (h *JWTHandler) IssueDefaultToken(subject string) string {
	method := jwt.SigningMethodHS256
	claims := claims{
		"sub": subject,
		"iat": time.Now().Unix(),
		"exp": time.Now().Unix() + 10,
	}
	secret, _ := h.keyFunc(nil)
	ss, err := jwt.NewWithClaims(method, claims).SignedString(secret)
	if err != nil {
		logger.Logger().Error().Err(err).Msg("issue default token failed")
	}
	return ss
}

type claims map[string]interface{}

func (claims) Valid() error {
	return nil
}
