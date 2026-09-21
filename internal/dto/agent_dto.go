package dto

import (
	"github.com/google/uuid"
)

type CreateAgentRequest struct {
	Name           string    `json:"name" binding:"required,min=3,max=255"`
	CreatorID      uuid.UUID `json:"creator_id" binding:"required"`
	Category       string    `json:"category" binding:"required,max=255"`
	AgentType      string    `json:"agent_type" binding:"required,max=255"`
	Runtime        string    `json:"runtime" binding:"required,max=255"`
	Description    *string   `json:"description"`
	ReadmeMarkdown *string   `json:"readme_markdown"`
}

type AgentResponse struct {
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
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}
