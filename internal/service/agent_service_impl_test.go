package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"marketplace/internal/model"
	"marketplace/internal/service"
)

// mockAgentRepository simule les appels au repository
type mockAgentRepository struct {
	getAllFn           func(ctx context.Context) ([]model.Agent, error)
	createFn           func(ctx context.Context, a *model.Agent) (*model.Agent, error)
	getByIDFn          func(ctx context.Context, id uuid.UUID) (*model.Agent, error)
	searchFn           func(ctx context.Context, query string) ([]model.Agent, error)
	getVersionByNameFn func(ctx context.Context, name, version string) (*model.Agent, *model.AgentVersion, error)
}

func (m *mockAgentRepository) GetAll(ctx context.Context) ([]model.Agent, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx)
	}
	return nil, nil
}

func (m *mockAgentRepository) Create(ctx context.Context, a *model.Agent) (*model.Agent, error) {
	if m.createFn != nil {
		return m.createFn(ctx, a)
	}
	return nil, nil
}

func (m *mockAgentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Agent, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockAgentRepository) Search(ctx context.Context, query string) ([]model.Agent, error) {
	if m.searchFn != nil {
		return m.searchFn(ctx, query)
	}
	return nil, nil
}

func (m *mockAgentRepository) GetVersionByName(ctx context.Context, name, version string) (*model.Agent, *model.AgentVersion, error) {
	if m.getVersionByNameFn != nil {
		return m.getVersionByNameFn(ctx, name, version)
	}
	return nil, nil, nil
}

func TestAgentService_SearchPackages(t *testing.T) {
	ctx := context.Background()

	t.Run("Succès - Recherche déléguée au repo", func(t *testing.T) {
		expectedAgents := []model.Agent{
			{ID: uuid.New(), Name: "odin-crawler"},
		}
		repo := &mockAgentRepository{
			searchFn: func(ctx context.Context, query string) ([]model.Agent, error) {
				assert.Equal(t, "crawler", query)
				return expectedAgents, nil
			},
		}

		svc := service.NewAgentService(repo)
		agents, err := svc.SearchPackages(ctx, "crawler")

		assert.NoError(t, err)
		assert.Equal(t, expectedAgents, agents)
	})

	t.Run("Erreur SQL propagée", func(t *testing.T) {
		repo := &mockAgentRepository{
			searchFn: func(ctx context.Context, query string) ([]model.Agent, error) {
				return nil, errors.New("db error")
			},
		}

		svc := service.NewAgentService(repo)
		agents, err := svc.SearchPackages(ctx, "crawler")

		assert.Error(t, err)
		assert.Nil(t, agents)
	})
}

func TestAgentService_GetPackageVersion(t *testing.T) {
	ctx := context.Background()
	agentID := uuid.New()
	mockAgent := &model.Agent{ID: agentID, Name: "odin-crawler"}
	mockVersion := &model.AgentVersion{ID: uuid.New(), AgentID: agentID, Version: "1.0.0", PublishedAt: time.Now()}

	t.Run("Succès - Agent et version trouvés", func(t *testing.T) {
		repo := &mockAgentRepository{
			getVersionByNameFn: func(ctx context.Context, name, version string) (*model.Agent, *model.AgentVersion, error) {
				assert.Equal(t, "odin-crawler", name)
				assert.Equal(t, "1.0.0", version)
				return mockAgent, mockVersion, nil
			},
		}

		svc := service.NewAgentService(repo)
		agent, ver, err := svc.GetPackageVersion(ctx, "odin-crawler", "1.0.0")

		assert.NoError(t, err)
		assert.Equal(t, mockAgent, agent)
		assert.Equal(t, mockVersion, ver)
	})

	t.Run("Erreur - ErrPackageNotFound quand l'agent est nil", func(t *testing.T) {
		repo := &mockAgentRepository{
			getVersionByNameFn: func(ctx context.Context, name, version string) (*model.Agent, *model.AgentVersion, error) {
				return nil, nil, nil
			},
		}

		svc := service.NewAgentService(repo)
		agent, ver, err := svc.GetPackageVersion(ctx, "unknown-agent", "1.0.0")

		assert.ErrorIs(t, err, service.ErrPackageNotFound)
		assert.Nil(t, agent)
		assert.Nil(t, ver)
	})

	t.Run("Erreur - ErrVersionNotFound quand la version est nil", func(t *testing.T) {
		repo := &mockAgentRepository{
			getVersionByNameFn: func(ctx context.Context, name, version string) (*model.Agent, *model.AgentVersion, error) {
				return mockAgent, nil, nil
			},
		}

		svc := service.NewAgentService(repo)
		agent, ver, err := svc.GetPackageVersion(ctx, "odin-crawler", "9.9.9")

		assert.ErrorIs(t, err, service.ErrVersionNotFound)
		assert.Nil(t, agent)
		assert.Nil(t, ver)
	})

	t.Run("Erreur SQL propagée", func(t *testing.T) {
		repo := &mockAgentRepository{
			getVersionByNameFn: func(ctx context.Context, name, version string) (*model.Agent, *model.AgentVersion, error) {
				return nil, nil, errors.New("db query failed")
			},
		}

		svc := service.NewAgentService(repo)
		agent, ver, err := svc.GetPackageVersion(ctx, "odin-crawler", "1.0.0")

		assert.Error(t, err)
		assert.Nil(t, agent)
		assert.Nil(t, ver)
	})
}
