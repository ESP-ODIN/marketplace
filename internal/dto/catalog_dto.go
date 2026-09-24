package dto

import (
	"time"

	"marketplace/internal/model"

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
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func ToAgentResponse(a model.Agent) AgentResponse {
	return AgentResponse{
		ID:             a.ID,
		Name:           a.Name,
		CreatorID:      a.CreatorID,
		Category:       a.Category,
		AgentType:      a.AgentType,
		Runtime:        a.Runtime,
		Description:    a.Description,
		ReadmeMarkdown: a.ReadmeMarkdown,
		IsOfficialPick: a.IsOfficialPick,
		DownloadsCount: a.DownloadsCount,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}
}

// ToAgentResponseList mappe un slice de model.Agent vers un slice de dto.AgentResponse
func ToAgentResponseList(agents []model.Agent) []AgentResponse {
	res := make([]AgentResponse, len(agents))
	for i, a := range agents {
		res[i] = ToAgentResponse(a)
	}
	return res
}

type ErrorResponse struct {
	Error   string `json:"error" example:"Données invalides"`
	Details string `json:"details,omitempty" example:"Le champ name est obligatoire"`
}

type SingleAgentResponse struct {
	Data model.Agent `json:"data"`
}

type AgentListResponse struct {
	Data []model.Agent `json:"data"`
}
