package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"marketplace/internal/dto"
	"marketplace/internal/handler"
	"marketplace/internal/model"
	"marketplace/internal/service"
)

type profileAgentService struct {
	getProfileFn func(ctx context.Context, creatorID uuid.UUID) ([]model.Agent, error)
}

func (m *profileAgentService) GetAllAgents(context.Context) ([]model.Agent, error) {
	return nil, nil
}
func (m *profileAgentService) CreateAgent(context.Context, dto.CreateAgentRequest) (*model.Agent, error) {
	return nil, nil
}
func (m *profileAgentService) GetAgentByID(context.Context, uuid.UUID) (*model.Agent, error) {
	return nil, nil
}
func (m *profileAgentService) GetProfile(ctx context.Context, creatorID uuid.UUID) ([]model.Agent, error) {
	return m.getProfileFn(ctx, creatorID)
}
func (m *profileAgentService) SearchPackages(context.Context, string) ([]model.Agent, error) {
	return nil, nil
}
func (m *profileAgentService) GetPackageVersion(context.Context, string, string) (*model.Agent, *model.AgentVersion, error) {
	return nil, nil, nil
}

func setupProfileRouter(svc service.AgentService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewAgentHandler(svc)
	r.GET("/api/v1/catalog/profiles/:id", h.GetProfile)
	return r
}

func TestAgentHandler_GetProfile(t *testing.T) {
	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	t.Run("400 - uuid invalide", func(t *testing.T) {
		r := setupProfileRouter(&profileAgentService{
			getProfileFn: func(context.Context, uuid.UUID) ([]model.Agent, error) {
				t.Fatal("le service ne doit pas être appelé")
				return nil, nil
			},
		})
		req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/profiles/pas-un-uuid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var body dto.ErrorResponse
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		assert.Equal(t, "l'identifiant fourni n'est pas un UUID valide", body.Error)
	})

	t.Run("200 - agents du créateur", func(t *testing.T) {
		agentID := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		r := setupProfileRouter(&profileAgentService{
			getProfileFn: func(_ context.Context, id uuid.UUID) ([]model.Agent, error) {
				assert.Equal(t, creatorID, id)
				return []model.Agent{{
					ID:        agentID,
					Name:      "odin-crawler",
					CreatorID: creatorID,
					Category:  "search",
					AgentType: "tool",
					Runtime:   "python",
					CreatedAt: time.Date(2026, 10, 22, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2026, 10, 22, 0, 0, 0, 0, time.UTC),
				}}, nil
			},
		})
		req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/profiles/"+creatorID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var body struct {
			Data []dto.AgentResponse `json:"data"`
		}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		assert.Len(t, body.Data, 1)
		assert.Equal(t, "odin-crawler", body.Data[0].Name)
		assert.Equal(t, agentID, body.Data[0].ID)
		assert.Equal(t, creatorID, body.Data[0].CreatorID)
	})

	t.Run("500 - erreur du service", func(t *testing.T) {
		r := setupProfileRouter(&profileAgentService{
			getProfileFn: func(context.Context, uuid.UUID) ([]model.Agent, error) {
				return nil, errors.New("db down")
			},
		})
		req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/profiles/"+creatorID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var body dto.ErrorResponse
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		assert.Equal(t, "erreur interne lors de la récupération du profil", body.Error)
	})
}
