package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"marketplace/internal/dto"
	"marketplace/internal/model"
)

var (
	ErrPackageNotFound = errors.New("package not found")
	ErrVersionNotFound = errors.New("version not found")
)

type AgentService interface {
	GetAgentByID(ctx context.Context, id uuid.UUID) (*model.Agent, error)
	GetAllAgents(ctx context.Context) ([]model.Agent, error)
	CreateAgent(ctx context.Context, req dto.CreateAgentRequest) (*model.Agent, error)
	SearchPackages(ctx context.Context, query string) ([]model.Agent, error)
	GetPackageVersion(ctx context.Context, name, version string) (*model.Agent, *model.AgentVersion, error)
}
