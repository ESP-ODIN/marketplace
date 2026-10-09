package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"marketplace/internal/repository"
)

func TestAgentRepository_Integration(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL non définie, skip du test d'intégration repository")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	assert.NoError(t, err)
	defer pool.Close()

	repo := repository.NewAgentRepository(pool)

	t.Run("Search", func(t *testing.T) {
		agents, err := repo.Search(ctx, "")
		assert.NoError(t, err)
		assert.NotNil(t, agents)
	})

	t.Run("GetVersionByName - Inexistant", func(t *testing.T) {
		agent, ver, err := repo.GetVersionByName(ctx, "non-existent-agent-xyz-404", "latest")
		assert.NoError(t, err)
		assert.Nil(t, agent)
		assert.Nil(t, ver)
	})
}
