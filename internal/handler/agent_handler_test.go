package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"marketplace/internal/dto"
	"marketplace/internal/model"
	"marketplace/internal/repository"
)

func TestAgentHandler_Create(t *testing.T) {
	validBody := `{
		"name": "odin-crawler",
		"category": "data",
		"agent_type": "tool",
		"runtime": "go"
	}`

	cases := []struct {
		name       string
		body       string
		svcErr     error
		wantStatus int
		wantError  string
	}{
		{"409 - nom déjà utilisé", validBody, repository.ErrAgentAlreadyExists, http.StatusConflict, "un agent avec ce nom existe déjà"},
		{"400 - créateur inconnu", validBody, repository.ErrCreatorNotFound, http.StatusBadRequest, "le créateur spécifié n'existe pas"},
		{"500 - erreur inattendue", validBody, errors.New("db down"), http.StatusInternalServerError, "erreur interne lors de la création de l'agent"},
		{"400 - JSON invalide sans détail interne", `{"name": 42}`, nil, http.StatusBadRequest, "données invalides"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockAgentService{
				createFn: func(ctx context.Context, input dto.CreateAgentRequest) (*model.Agent, error) {
					return nil, tc.svcErr
				},
			}

			r, signForUser := setupCatalogAuthRouter(t, mockSvc)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/catalog/agents", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+signForUser(uuid.New()))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)

			var res dto.ErrorResponse
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.Equal(t, tc.wantError, res.Error)
			assert.Empty(t, res.Details, "aucun détail interne ne doit être exposé")
		})
	}

	t.Run("201 - agent créé", func(t *testing.T) {
		userID := uuid.New()
		mockSvc := &mockAgentService{
			createFn: func(ctx context.Context, input dto.CreateAgentRequest) (*model.Agent, error) {
				return &model.Agent{ID: uuid.New(), Name: input.Name, CreatorID: input.CreatorID}, nil
			},
		}

		r, signForUser := setupCatalogAuthRouter(t, mockSvc)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/catalog/agents", strings.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+signForUser(userID))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var res dto.SingleAgentResponse
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Equal(t, "odin-crawler", res.Data.Name)
		assert.Equal(t, userID, res.Data.CreatorID)
	})
}
