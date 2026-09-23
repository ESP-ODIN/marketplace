package model

import (
	"time"

	"github.com/google/uuid"
)

type Agent struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	CreatorID      uuid.UUID `json:"creator_id"`
	Category       string    `json:"category"`
	AgentType      string    `json:"agent_type"`
	Runtime        string    `json:"runtime"`
	Description    *string   `json:"description,omitempty"`
	ReadmeMarkdown *string   `json:"readme_markdown,omitempty"`
	IsOfficialPick bool      `json:"is_official_pick"`
	DownloadsCount int       `json:"downloads_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
