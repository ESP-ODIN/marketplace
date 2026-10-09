package dto

import "marketplace/internal/model"

// PackageSearchResponse représente la liste de packages retournée par la recherche
type PackageSearchResponse struct {
	Count int           `json:"count"`
	Data  []model.Agent `json:"data"`
}

// PackageDetailResponse représente un package et les métadonnées de sa version
type PackageDetailResponse struct {
	Agent   model.Agent        `json:"agent"`
	Version model.AgentVersion `json:"version"`
}

type CLIErrorResponse struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}
