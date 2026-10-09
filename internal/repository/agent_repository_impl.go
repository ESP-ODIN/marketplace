package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace/internal/model"
)

type agentRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewAgentRepository(db *pgxpool.Pool) AgentRepository {
	return &agentRepositoryImpl{db: db}
}

func (r *agentRepositoryImpl) GetAll(ctx context.Context) ([]model.Agent, error) {
	query := `
		SELECT id, name, creator_id, category, agent_type, runtime, 
		       description, readme_markdown, is_official_pick, downloads_count, 
		       created_at, updated_at 
		FROM agent
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
			&a.ID, &a.Name, &a.CreatorID, &a.Category, &a.AgentType, &a.Runtime,
			&a.Description, &a.ReadmeMarkdown, &a.IsOfficialPick, &a.DownloadsCount,
			&a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		agents = append(agents, a)
	}

	return agents, nil
}

func (r *agentRepositoryImpl) Create(ctx context.Context, a *model.Agent) (*model.Agent, error) {
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

func (r *agentRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Agent, error) {
	query := `
		SELECT id, creator_id, name, description, readme_markdown, category, agent_type, runtime, created_at, updated_at
		FROM agent
		WHERE id = $1;
	`

	var a model.Agent
	err := r.db.QueryRow(ctx, query, id).Scan(
		&a.ID,
		&a.CreatorID,
		&a.Name,
		&a.Description,
		&a.ReadmeMarkdown,
		&a.Category,
		&a.AgentType,
		&a.Runtime,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("agentRepository.GetByID: %w", err)
	}

	return &a, nil
}

func (r *agentRepositoryImpl) Search(ctx context.Context, q string) ([]model.Agent, error) {
	searchTerm := "%" + q + "%"
	querySQL := `
		SELECT id, name, creator_id, category, agent_type, runtime, 
		       description, readme_markdown, is_official_pick, downloads_count, 
		       created_at, updated_at
		FROM agent
		WHERE name ILIKE $1 OR description ILIKE $1
		ORDER BY downloads_count DESC, created_at DESC;
	`

	rows, err := r.db.Query(ctx, querySQL, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("agentRepository.Search: %w", err)
	}
	defer rows.Close()

	agents := make([]model.Agent, 0)
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
			return nil, fmt.Errorf("agentRepository.Search scan: %w", err)
		}
		agents = append(agents, a)
	}

	return agents, nil
}

func (r *agentRepositoryImpl) GetVersionByName(ctx context.Context, name string, version string) (*model.Agent, *model.AgentVersion, error) {
	// 1. Vérifier si l'agent existe
	agentQuery := `
		SELECT id, creator_id, name, description, readme_markdown, category, agent_type, runtime, is_official_pick, downloads_count, created_at, updated_at
		FROM agent
		WHERE name = $1;
	`
	var a model.Agent
	err := r.db.QueryRow(ctx, agentQuery, name).Scan(
		&a.ID, &a.CreatorID, &a.Name, &a.Description, &a.ReadmeMarkdown,
		&a.Category, &a.AgentType, &a.Runtime, &a.IsOfficialPick,
		&a.DownloadsCount, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, nil // Agent inconnu
		}
		return nil, nil, fmt.Errorf("agentRepository.GetVersionByName (agent): %w", err)
	}

	// 2. Chercher la version demandée (ou latest)
	var versionQuery string
	var args []any

	if version == "" || version == "latest" {
		versionQuery = `
			SELECT id, agent_id, version, sha256, manifest_json, changelog, is_stable, security_run_id, published_at, repo_link
			FROM agent_version
			WHERE agent_id = $1
			ORDER BY is_stable DESC, published_at DESC
			LIMIT 1;
		`
		args = []any{a.ID}
	} else {
		versionQuery = `
			SELECT id, agent_id, version, sha256, manifest_json, changelog, is_stable, security_run_id, published_at, repo_link
			FROM agent_version
			WHERE agent_id = $1 AND version = $2;
		`
		args = []any{a.ID, version}
	}

	var v model.AgentVersion
	err = r.db.QueryRow(ctx, versionQuery, args...).Scan(
		&v.ID, &v.AgentID, &v.Version, &v.Sha256, &v.ManifestJSON,
		&v.Changelog, &v.IsStable, &v.SecurityRunID, &v.PublishedAt, &v.RepoLink,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &a, nil, nil // Agent trouvé mais version introuvable
		}
		return nil, nil, fmt.Errorf("agentRepository.GetVersionByName (version): %w", err)
	}

	return &a, &v, nil
}
