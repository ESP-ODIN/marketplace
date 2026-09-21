package repository

import (
	"context"

	"marketplace/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AgentRepository interface {
	GetAll(ctx context.Context) ([]model.Agent, error)
	Create(ctx context.Context, agent *model.Agent) (*model.Agent, error)
}

type agentRepository struct {
	db *pgxpool.Pool
}

func NewAgentRepository(db *pgxpool.Pool) AgentRepository {
	return &agentRepository{db: db}
}

func (r *agentRepository) GetAll(ctx context.Context) ([]model.Agent, error) {
	query := `
		SELECT 
			id, name, creator_id, category, agent_type, runtime, 
			description, readme_markdown, is_official_pick, downloads_count, 
			created_at, updated_at
		FROM agent
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []model.Agent
	for rows.Next() {
		var a model.Agent
		err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.CreatorID,
			&a.Category,
			&a.AgentType,
			&a.Runtime,
			&a.Description,
			&a.ReadmeMarkdown,
			&a.IsOfficialPick,
			&a.DownloadsCount,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		agents = append(agents, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return agents, nil
}

func (r *agentRepository) Create(ctx context.Context, a *model.Agent) (*model.Agent, error) {
	query := `
		INSERT INTO agent (
			name, creator_id, category, agent_type, runtime, description, readme_markdown
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, is_official_pick, downloads_count, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		a.Name,
		a.CreatorID,
		a.Category,
		a.AgentType,
		a.Runtime,
		a.Description,
		a.ReadmeMarkdown,
	).Scan(
		&a.ID,
		&a.IsOfficialPick,
		&a.DownloadsCount,
		&a.CreatedAt,
		&a.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return a, nil
}
