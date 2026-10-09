package repository

import (
	"context"

	"github.com/google/uuid"

	"marketplace/internal/model"
)

type AgentRepository interface {
	GetAll(ctx context.Context) ([]model.Agent, error)
	Create(ctx context.Context, a *model.Agent) (*model.Agent, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Agent, error)
	Search(ctx context.Context, query string) ([]model.Agent, error)
	GetVersionByName(ctx context.Context, name string, version string) (*model.Agent, *model.AgentVersion, error)
}
