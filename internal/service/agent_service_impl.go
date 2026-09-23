package service

import (
	"context"

	"marketplace/internal/dto"
	"marketplace/internal/model"
	"marketplace/internal/repository"
)

type agentServiceImpl struct {
	repo repository.AgentRepository
}

func NewAgentService(repo repository.AgentRepository) AgentService {
	return &agentServiceImpl{repo: repo}
}

func (s *agentServiceImpl) GetAllAgents(ctx context.Context) ([]model.Agent, error) {
	return s.repo.GetAll(ctx)
}

func (s *agentServiceImpl) CreateAgent(ctx context.Context, req dto.CreateAgentRequest) (*model.Agent, error) {
	agent := &model.Agent{
		Name:           req.Name,
		CreatorID:      req.CreatorID,
		Category:       req.Category,
		AgentType:      req.AgentType,
		Runtime:        req.Runtime,
		Description:    req.Description,
		ReadmeMarkdown: req.ReadmeMarkdown,
	}

	return s.repo.Create(ctx, agent)
}
