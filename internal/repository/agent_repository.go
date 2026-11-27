package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"marketplace/internal/model"
)

var (
	// ErrAgentAlreadyExists : un agent avec ce nom existe déjà (contrainte unique, code PostgreSQL 23505).
	ErrAgentAlreadyExists = errors.New("agent already exists")
	// ErrCreatorNotFound : le creator_id ne correspond à aucun profil (clé étrangère, code PostgreSQL 23503).
	ErrCreatorNotFound = errors.New("creator not found")
)

type AgentRepository interface {
	GetAll(ctx context.Context) ([]model.Agent, error)
	Create(ctx context.Context, a *model.Agent) (*model.Agent, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Agent, error)
	ListByCreatorID(ctx context.Context, creatorID uuid.UUID) ([]model.Agent, error)
	Search(ctx context.Context, query string) ([]model.Agent, error)
	GetVersionByName(ctx context.Context, name string, version string) (*model.Agent, *model.AgentVersion, error)
}
