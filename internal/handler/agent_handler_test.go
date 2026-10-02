package handler

import (
	"bytes"
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
	"marketplace/internal/model"
	"marketplace/internal/service"
)

type mockCatalogAgentService struct {
	getAllAgentsFn func(ctx context.Context) ([]model.Agent, error)
	getAgentByIDFn func(ctx context.Context, id uuid.UUID) (*model.Agent, error)
	createAgentFn  func(ctx context.Context, input dto.CreateAgentRequest) (*model.Agent, error)
}

func (m *mockCatalogAgentService) GetAllAgents(ctx context.Context) ([]model.Agent, error) {
	if m.getAllAgentsFn != nil {
		return m.getAllAgentsFn(ctx)
	}
	return nil, nil
}

func (m *mockCatalogAgentService) GetAgentByID(ctx context.Context, id uuid.UUID) (*model.Agent, error) {
	if m.getAgentByIDFn != nil {
		return m.getAgentByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockCatalogAgentService) CreateAgent(ctx context.Context, input dto.CreateAgentRequest) (*model.Agent, error) {
	if m.createAgentFn != nil {
		return m.createAgentFn(ctx, input)
	}
	return nil, nil
}

func (m *mockCatalogAgentService) SearchPackages(ctx context.Context, query string) ([]model.Agent, error) {
	return nil, nil
}

func (m *mockCatalogAgentService) GetPackageVersion(ctx context.Context, name, version string) (*model.Agent, *model.AgentVersion, error) {
	return nil, nil, nil
}

func setupAgentRouter(mock service.AgentService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewAgentHandler(mock)

	catalog := r.Group("/catalog/agents")
	{
		catalog.GET("", h.GetAll)
		catalog.GET("/:id", h.GetByID)
		catalog.POST("", h.Create)
	}

	return r
}

func TestAgentHandler_GetAll(t *testing.T) {
	t.Run("200 - retourne la liste des agents", func(t *testing.T) {
		agentID := uuid.New()
		creatorID := uuid.New()
		mockAgents := []model.Agent{
			{
				ID:        agentID,
				Name:      "odin-agent",
				CreatorID: creatorID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		mockSvc := &mockCatalogAgentService{
			getAllAgentsFn: func(ctx context.Context) ([]model.Agent, error) {
				return mockAgents, nil
			},
		}

		r := setupAgentRouter(mockSvc)
		req, _ := http.NewRequest(http.MethodGet, "/catalog/agents", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var res map[string][]dto.AgentResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Len(t, res["data"], 1)
		assert.Equal(t, "odin-agent", res["data"][0].Name)
	})

	t.Run("500 - erreur lors de la recuperation", func(t *testing.T) {
		mockSvc := &mockCatalogAgentService{
			getAllAgentsFn: func(ctx context.Context) ([]model.Agent, error) {
				return nil, errors.New("db error")
			},
		}

		r := setupAgentRouter(mockSvc)
		req, _ := http.NewRequest(http.MethodGet, "/catalog/agents", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAgentHandler_GetByID(t *testing.T) {
	agentID := uuid.New()

	t.Run("200 - agent trouve", func(t *testing.T) {
		mockAgent := &model.Agent{
			ID:        agentID,
			Name:      "odin-agent",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockSvc := &mockCatalogAgentService{
			getAgentByIDFn: func(ctx context.Context, id uuid.UUID) (*model.Agent, error) {
				assert.Equal(t, agentID, id)
				return mockAgent, nil
			},
		}

		r := setupAgentRouter(mockSvc)
		req, _ := http.NewRequest(http.MethodGet, "/catalog/agents/"+agentID.String(), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var res dto.SingleAgentResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, "odin-agent", res.Data.Name)
	})

	t.Run("400 - UUID invalide", func(t *testing.T) {
		r := setupAgentRouter(&mockCatalogAgentService{})
		req, _ := http.NewRequest(http.MethodGet, "/catalog/agents/invalid-uuid-string", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var res dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Contains(t, res.Error, "UUID valide")
	})

	t.Run("404 - agent introuvable", func(t *testing.T) {
		mockSvc := &mockCatalogAgentService{
			getAgentByIDFn: func(ctx context.Context, id uuid.UUID) (*model.Agent, error) {
				return nil, nil
			},
		}

		r := setupAgentRouter(mockSvc)
		req, _ := http.NewRequest(http.MethodGet, "/catalog/agents/"+agentID.String(), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("500 - erreur interne", func(t *testing.T) {
		mockSvc := &mockCatalogAgentService{
			getAgentByIDFn: func(ctx context.Context, id uuid.UUID) (*model.Agent, error) {
				return nil, errors.New("db timeout")
			},
		}

		r := setupAgentRouter(mockSvc)
		req, _ := http.NewRequest(http.MethodGet, "/catalog/agents/"+agentID.String(), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAgentHandler_Create(t *testing.T) {
	creatorID := uuid.New()
	desc := "description de test"
	readme := "# Readme"

	validPayload := dto.CreateAgentRequest{
		CreatorID:      creatorID,
		Name:           "new-agent",
		Description:    &desc,
		ReadmeMarkdown: &readme,
		Category:       "scraper",
		AgentType:      "agent",
		Runtime:        "python",
	}

	t.Run("201 - creation reussie", func(t *testing.T) {
		created := &model.Agent{
			ID:             uuid.New(),
			CreatorID:      creatorID,
			Name:           validPayload.Name,
			Description:    validPayload.Description,
			ReadmeMarkdown: validPayload.ReadmeMarkdown,
			Category:       validPayload.Category,
			AgentType:      validPayload.AgentType,
			Runtime:        validPayload.Runtime,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		mockSvc := &mockCatalogAgentService{
			createAgentFn: func(ctx context.Context, input dto.CreateAgentRequest) (*model.Agent, error) {
				assert.Equal(t, "new-agent", input.Name)
				return created, nil
			},
		}

		body, _ := json.Marshal(validPayload)
		r := setupAgentRouter(mockSvc)
		req, _ := http.NewRequest(http.MethodPost, "/catalog/agents", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var res map[string]dto.AgentResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, "new-agent", res["data"].Name)
	})

	t.Run("400 - body JSON invalide", func(t *testing.T) {
		r := setupAgentRouter(&mockCatalogAgentService{})
		req, _ := http.NewRequest(http.MethodPost, "/catalog/agents", bytes.NewBufferString("{ invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("500 - echec service", func(t *testing.T) {
		mockSvc := &mockCatalogAgentService{
			createAgentFn: func(ctx context.Context, input dto.CreateAgentRequest) (*model.Agent, error) {
				return nil, errors.New("insert failed")
			},
		}

		body, _ := json.Marshal(validPayload)
		r := setupAgentRouter(mockSvc)
		req, _ := http.NewRequest(http.MethodPost, "/catalog/agents", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
