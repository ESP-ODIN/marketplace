package model

import (
	"time"

	"github.com/google/uuid"
)

type AgentVersion struct {
	ID            uuid.UUID  `json:"id"`
	AgentID       uuid.UUID  `json:"agent_id"`
	Version       string     `json:"version"`
	Sha256        string     `json:"sha256"`
	ManifestJSON  *string    `json:"manifest_json,omitempty"`
	Changelog     *string    `json:"changelog,omitempty"`
	IsStable      bool       `json:"is_stable"`
	SecurityRunID *uuid.UUID `json:"security_run_id,omitempty"`
	PublishedAt   time.Time  `json:"published_at"`
	RepoLink      *string    `json:"repo_link,omitempty"`
}
