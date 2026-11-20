package middleware_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ESP-ODIN/authkit-go"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"marketplace/internal/middleware"
)

type jwksEnv struct {
	privateKey *rsa.PrivateKey
	kid        string
	issuer     string
	audience   string
	server     *httptest.Server
	auth       *authkit.Authenticator
}

func newTestJWKSEnv(t *testing.T) *jwksEnv {
	t.Helper()

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	kid := "test-key-1"
	jwksPayload := map[string]any{
		"keys": []map[string]any{
			{
				"kty": "RSA",
				"alg": "RS256",
				"use": "sig",
				"kid": kid,
				"n":   base64.RawURLEncoding.EncodeToString(privKey.PublicKey.N.Bytes()),
				"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(privKey.PublicKey.E)).Bytes()),
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jwksPayload)
	}))
	t.Cleanup(srv.Close)

	issuer := "https://auth.odin.test"
	audience := "odin-api"

	authenticator, err := authkit.New(authkit.Config{
		JWKSURL:    srv.URL,
		Issuer:     issuer,
		Audience:   audience,
		HTTPClient: srv.Client(),
	})
	require.NoError(t, err)

	return &jwksEnv{
		privateKey: privKey,
		kid:        kid,
		issuer:     issuer,
		audience:   audience,
		server:     srv,
		auth:       authenticator,
	}
}

func (e *jwksEnv) signToken(t *testing.T, key *rsa.PrivateKey, claims authkit.Claims) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tok.Header["kid"] = e.kid
	signed, err := tok.SignedString(key)
	require.NoError(t, err)
	return signed
}

func TestRequireAuth_Middleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	env := newTestJWKSEnv(t)

	setupEngine := func(onSuccess func(c *gin.Context)) *gin.Engine {
		r := gin.New()
		r.POST("/protected", middleware.RequireAuth(env.auth), func(c *gin.Context) {
			if onSuccess != nil {
				onSuccess(c)
			}
			c.Status(http.StatusOK)
		})
		return r
	}

	t.Run("401 - en-tête Authorization manquant", func(t *testing.T) {
		called := false
		r := setupEngine(func(*gin.Context) { called = true })

		req := httptest.NewRequest(http.MethodPost, "/protected", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.False(t, called)
	})

	t.Run("401 - jeton signé par une autre clé RSA", func(t *testing.T) {
		otherKey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)

		token := env.signToken(t, otherKey, authkit.Claims{
			ID:    uuid.New(),
			Email: "attacker@odin.test",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    env.issuer,
				Audience:  jwt.ClaimStrings{env.audience},
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		})

		called := false
		r := setupEngine(func(*gin.Context) { called = true })

		req := httptest.NewRequest(http.MethodPost, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.False(t, called)
	})

	t.Run("401 - jeton expiré", func(t *testing.T) {
		token := env.signToken(t, env.privateKey, authkit.Claims{
			ID:    uuid.New(),
			Email: "user@odin.test",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    env.issuer,
				Audience:  jwt.ClaimStrings{env.audience},
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Minute)),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			},
		})

		r := setupEngine(nil)
		req := httptest.NewRequest(http.MethodPost, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("200 - jeton valide injecte Identity dans le contexte", func(t *testing.T) {
		userID := uuid.New()
		token := env.signToken(t, env.privateKey, authkit.Claims{
			ID:    userID,
			Email: "creator@odin.test",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    env.issuer,
				Audience:  jwt.ClaimStrings{env.audience},
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		})

		var gotIdentity authkit.Identity
		var hasIdentity bool
		r := setupEngine(func(c *gin.Context) {
			gotIdentity, hasIdentity = authkit.IdentityFromContext(c.Request.Context())
		})

		req := httptest.NewRequest(http.MethodPost, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, hasIdentity)
		assert.Equal(t, userID, gotIdentity.ID)
		assert.Equal(t, "creator@odin.test", gotIdentity.Email)
	})
}
