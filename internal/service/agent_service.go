package service

import (
	"context"
	"uuid"

	"marketplace/internal/dto"
	"marketplace/internal/model"
)

type AgentService interface {
	GetAgentByID(ctx context.Context, id uuid.UUID) (*model.Agent, error)
	GetAllAgents(ctx context.Context) ([]model.Agent, error)
	CreateAgent(ctx context.Context, req dto.CreateAgentRequest) (*model.Agent, error)
}
