package repository_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"marketplace/internal/model"
	"marketplace/internal/repository"
)

type fakeRowScanner struct {
	values []any
}

func (f fakeRowScanner) Scan(dest ...any) error {
	for i := range dest {
		switch d := dest[i].(type) {
		case *uuid.UUID:
			*d = f.values[i].(uuid.UUID)
		case *string:
			*d = f.values[i].(string)
		case **string:
			if f.values[i] != nil {
				s := f.values[i].(string)
				*d = &s
			}
		case *bool:
			*d = f.values[i].(bool)
		case *int:
			*d = f.values[i].(int)
		case *time.Time:
			*d = f.values[i].(time.Time)
		}
	}
	return nil
}

func TestScanAgent_ReadsOfficialPickAndDownloadsCount(t *testing.T) {
	id := uuid.New()
	creatorID := uuid.New()
	now := time.Now().UTC()
	desc := "Description de test"

	scanner := fakeRowScanner{
		values: []any{
			id,
			"odin-reviewer",
			creatorID,
			"devtools",
			"workflow",
			"typescript@5.4",
			desc,
			nil,
			true,
			4250,
			now,
			now,
		},
	}

	got, err := repository.ScanAgentForTest(scanner)
	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, "odin-reviewer", got.Name)
	assert.True(t, got.IsOfficialPick)
	assert.Equal(t, 4250, got.DownloadsCount)
}

func TestAgentRepository_Integration(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL non définie, skip du test d'intégration repository")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
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

	t.Run("GetByID - is_official_pick et downloads_count", func(t *testing.T) {
		var creatorID uuid.UUID
		err := pool.QueryRow(ctx, `SELECT id FROM profile LIMIT 1`).Scan(&creatorID)
		if err != nil {
			t.Skip("Aucun profil en base pour satisfaire agent_creator_id_fkey")
		}

		created, err := repo.Create(ctx, &model.Agent{
			Name:      "test-getbyid-" + uuid.NewString()[:8],
			CreatorID: creatorID,
			Category:  "test",
			AgentType: "tool",
			Runtime:   "go",
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			_, _ = pool.Exec(context.Background(), `DELETE FROM agent WHERE id = $1`, created.ID)
		})

		_, err = pool.Exec(ctx,
			`UPDATE agent SET is_official_pick = true, downloads_count = 1337 WHERE id = $1`,
			created.ID,
		)
		require.NoError(t, err)

		fetched, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		require.NotNil(t, fetched)
		assert.True(t, fetched.IsOfficialPick)
		assert.Equal(t, 1337, fetched.DownloadsCount)
	})
}
