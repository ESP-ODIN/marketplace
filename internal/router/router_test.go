package router_test

import (
	"context"
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
	"marketplace/internal/router"
)

type mockService struct{}

func (m *mockService) GetAllAgents(ctx context.Context) ([]model.Agent, error) {
	return []model.Agent{}, nil
}
func (m *mockService) CreateAgent(ctx context.Context, input dto.CreateAgentRequest) (*model.Agent, error) {
	return nil, nil
}
func (m *mockService) GetAgentByID(ctx context.Context, id uuid.UUID) (*model.Agent, error) {
	return nil, nil
}
func (m *mockService) SearchPackages(ctx context.Context, query string) ([]model.Agent, error) {
	return []model.Agent{}, nil
}
func (m *mockService) GetPackageVersion(ctx context.Context, name, version string) (*model.Agent, *model.AgentVersion, error) {
	agentID := uuid.New()
	return &model.Agent{
		ID:   agentID,
		Name: name,
	}, &model.AgentVersion{
		ID:          uuid.New(),
		AgentID:     agentID,
		Version:     version,
		PublishedAt: time.Now(),
	}, nil
}

func TestSetup_RoutesRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockService{}
	cfg := router.Config{
		HealthHandler:  handler.NewHealthHandler(nil),
		AgentHandler:   handler.NewAgentHandler(mockSvc),
		PackageHandler: handler.NewPackageHandler(mockSvc),
	}

	r := router.Setup(cfg)

	// 1. Vérification que la route wildcard Swagger est bien enregistrée
	routes := r.Routes()
	hasSwagger := false
	for _, route := range routes {
		if route.Path == "/swagger/*any" {
			hasSwagger = true
			break
		}
	}
	assert.True(t, hasSwagger, "la route Swagger /swagger/*any doit être enregistrée")

	// 2. Tests des requêtes HTTP sur l'API
	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{
			name:           "Health Live",
			method:         http.MethodGet,
			url:            "/health/live",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Catalog - List Agents",
			method:         http.MethodGet,
			url:            "/api/v1/catalog/agents",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Catalog - Create Agent Validation Error",
			method:         http.MethodPost,
			url:            "/api/v1/catalog/agents",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Packages - Search",
			method:         http.MethodGet,
			url:            "/api/v1/packages/search?q=test",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Packages - Get Version",
			method:         http.MethodGet,
			url:            "/api/v1/packages/test-agent/latest",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Route inexistante - 404",
			method:         http.MethodGet,
			url:            "/api/v1/unknown-endpoint",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
