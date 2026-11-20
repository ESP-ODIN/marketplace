package handler_test

import (
	"bytes"
	"context"
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

	"marketplace/internal/dto"
	"marketplace/internal/handler"
	"marketplace/internal/middleware"
	"marketplace/internal/model"
	"marketplace/internal/service"
)

type createAgentMockService struct {
	mockAgentService
	createAgentFn func(ctx context.Context, input dto.CreateAgentRequest) (*model.Agent, error)
}

func (m *createAgentMockService) CreateAgent(ctx context.Context, input dto.CreateAgentRequest) (*model.Agent, error) {
	if m.createAgentFn != nil {
		return m.createAgentFn(ctx, input)
	}
	return nil, nil
}

func setupCatalogAuthRouter(t *testing.T, svc service.AgentService) (*gin.Engine, func(userID uuid.UUID) string) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	kid := "catalog-test-kid"
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

	r := gin.New()
	h := handler.NewAgentHandler(svc)
	r.POST("/api/v1/catalog/agents", middleware.RequireAuth(authenticator), h.Create)

	signForUser := func(userID uuid.UUID) string {
		tok := jwt.NewWithClaims(jwt.SigningMethodRS256, authkit.Claims{
			ID:    userID,
			Email: "dev@odin.test",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    issuer,
				Audience:  jwt.ClaimStrings{audience},
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		})
		tok.Header["kid"] = kid
		signed, err := tok.SignedString(privKey)
		require.NoError(t, err)
		return signed
	}

	return r, signForUser
}

func TestAgentHandler_Create_WithAuth(t *testing.T) {
	t.Run("401 - sans jeton Bearer", func(t *testing.T) {
		r, _ := setupCatalogAuthRouter(t, &createAgentMockService{})

		body := `{"name":"odin-agent","category":"dev","agent_type":"tool","runtime":"go"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/catalog/agents", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("201 - CreatorID forcé depuis le JWT même si fourni dans le JSON", func(t *testing.T) {
		jwtUserID := uuid.New()
		spoofedCreatorID := uuid.New()
		createdAgentID := uuid.New()

		svc := &createAgentMockService{
			createAgentFn: func(_ context.Context, input dto.CreateAgentRequest) (*model.Agent, error) {
				assert.Equal(t, jwtUserID, input.CreatorID)
				assert.NotEqual(t, spoofedCreatorID, input.CreatorID)
				assert.Equal(t, "odin-agent", input.Name)

				return &model.Agent{
					ID:        createdAgentID,
					Name:      input.Name,
					CreatorID: input.CreatorID,
					Category:  input.Category,
					AgentType: input.AgentType,
					Runtime:   input.Runtime,
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
				}, nil
			},
		}

		r, signForUser := setupCatalogAuthRouter(t, svc)

		payload := map[string]any{
			"name":       "odin-agent",
			"creator_id": spoofedCreatorID.String(),
			"category":   "dev",
			"agent_type": "tool",
			"runtime":    "go",
		}
		rawBody, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/catalog/agents", bytes.NewReader(rawBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+signForUser(jwtUserID))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp struct {
			Data dto.AgentResponse `json:"data"`
		}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, createdAgentID, resp.Data.ID)
		assert.Equal(t, jwtUserID, resp.Data.CreatorID)
	})
}
