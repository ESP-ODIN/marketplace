package service

import (
	"context"

	"marketplace/internal/dto"
	"marketplace/internal/model"
)

type AgentService interface {
	GetAllAgents(ctx context.Context) ([]model.Agent, error)
	CreateAgent(ctx context.Context, req dto.CreateAgentRequest) (*model.Agent, error)
}
