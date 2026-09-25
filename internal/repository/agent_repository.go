package repository

import (
	"context"
	"uuid"

	"marketplace/internal/model"
)

type AgentRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Agent, error)
	GetAll(ctx context.Context) ([]model.Agent, error)
	Create(ctx context.Context, agent *model.Agent) (*model.Agent, error)
}
