package repository

import (
	"context"

	"marketplace/internal/model"
)

type AgentRepository interface {
	GetAll(ctx context.Context) ([]model.Agent, error)
	Create(ctx context.Context, agent *model.Agent) (*model.Agent, error)
}
