package handler_test

import (
	"context"
	"encoding/json"
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

// mockAgentService simule les réponses métier du service pour les tests HTTP
type mockAgentService struct {
	searchFn            func(ctx context.Context, query string) ([]model.Agent, error)
	getPackageVersionFn func(ctx context.Context, name, version string) (*model.Agent, *model.AgentVersion, error)
}

func (m *mockAgentService) GetAllAgents(ctx context.Context) ([]model.Agent, error) {
	return nil, nil
}
func (m *mockAgentService) CreateAgent(ctx context.Context, input dto.CreateAgentRequest) (*model.Agent, error) {
	return nil, nil
}
func (m *mockAgentService) GetAgentByID(ctx context.Context, id uuid.UUID) (*model.Agent, error) {
	return nil, nil
}
func (m *mockAgentService) SearchPackages(ctx context.Context, query string) ([]model.Agent, error) {
	if m.searchFn != nil {
		return m.searchFn(ctx, query)
	}
	return nil, nil
}
func (m *mockAgentService) GetPackageVersion(ctx context.Context, name, version string) (*model.Agent, *model.AgentVersion, error) {
	if m.getPackageVersionFn != nil {
		return m.getPackageVersionFn(ctx, name, version)
	}
	return nil, nil, nil
}

// setupRouter configure un moteur Gin de test isolé
func setupRouter(mock service.AgentService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewPackageHandler(mock)

	packages := r.Group("/api/v1/packages")
	{
		packages.GET("/search", h.Search)
		packages.GET("/:name/:version", h.GetByNameAndVersion)
	}

	return r
}

func TestPackageHandler_Search(t *testing.T) {
	t.Run("200 - retourne les résultats trouvés", func(t *testing.T) {
		mockSvc := &mockAgentService{
			searchFn: func(ctx context.Context, query string) ([]model.Agent, error) {
				assert.Equal(t, "crawler", query)
				return []model.Agent{
					{
						ID:   uuid.New(),
						Name: "odin-crawler",
					},
				}, nil
			},
		}

		r := setupRouter(mockSvc)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/packages/search?q=crawler", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var res dto.PackageSearchResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, 1, res.Count)
		assert.Equal(t, "odin-crawler", res.Data[0].Name)
	})

	t.Run("200 - recherche vide retourne liste vide", func(t *testing.T) {
		mockSvc := &mockAgentService{
			searchFn: func(ctx context.Context, query string) ([]model.Agent, error) {
				assert.Equal(t, "", query)
				return []model.Agent{}, nil
			},
		}

		r := setupRouter(mockSvc)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/packages/search", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var res dto.PackageSearchResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, 0, res.Count)
		assert.Empty(t, res.Data)
	})
}

func TestPackageHandler_GetByNameAndVersion(t *testing.T) {
	agentID := uuid.New()
	sampleAgent := &model.Agent{
		ID:   agentID,
		Name: "odin-crawler",
	}
	sampleVersion := &model.AgentVersion{
		ID:          uuid.New(),
		AgentID:     agentID,
		Version:     "1.0.0",
		Sha256:      "a3f5c78...",
		PublishedAt: time.Now(),
	}

	t.Run("200 - package et version trouvés avec 'latest'", func(t *testing.T) {
		mockSvc := &mockAgentService{
			getPackageVersionFn: func(ctx context.Context, name, version string) (*model.Agent, *model.AgentVersion, error) {
				assert.Equal(t, "odin-crawler", name)
				assert.Equal(t, "latest", version)
				return sampleAgent, sampleVersion, nil
			},
		}

		r := setupRouter(mockSvc)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/packages/odin-crawler/latest", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var res dto.PackageDetailResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, "odin-crawler", res.Agent.Name)
		assert.Equal(t, "1.0.0", res.Version.Version)
	})

	t.Run("404 - package introuvable (PACKAGE_NOT_FOUND)", func(t *testing.T) {
		mockSvc := &mockAgentService{
			getPackageVersionFn: func(ctx context.Context, name, version string) (*model.Agent, *model.AgentVersion, error) {
				return nil, nil, service.ErrPackageNotFound
			},
		}

		r := setupRouter(mockSvc)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/packages/unknown-pkg/1.0.0", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var res dto.CLIErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, "PACKAGE_NOT_FOUND", res.Code)
		assert.Equal(t, "unknown-pkg", res.Details["name"])
	})

	t.Run("404 - version introuvable pour ce package (VERSION_NOT_FOUND)", func(t *testing.T) {
		mockSvc := &mockAgentService{
			getPackageVersionFn: func(ctx context.Context, name, version string) (*model.Agent, *model.AgentVersion, error) {
				return nil, nil, service.ErrVersionNotFound
			},
		}

		r := setupRouter(mockSvc)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/packages/odin-crawler/9.9.9", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var res dto.CLIErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, "VERSION_NOT_FOUND", res.Code)
		assert.Equal(t, "9.9.9", res.Details["version"])
	})
}
