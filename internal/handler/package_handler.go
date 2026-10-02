package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"marketplace/internal/dto"
	"marketplace/internal/service"
)

type PackageHandler struct {
	service service.AgentService
}

func NewPackageHandler(s service.AgentService) *PackageHandler {
	return &PackageHandler{service: s}
}

// Search godoc
// @Summary      Rechercher des packages
// @Description  Recherche un package par mot-clé (retourne une liste vide si aucun mot-clé fourni)
// @Tags         packages
// @Produce      json
// @Param        q    query     string  false  "Terme de recherche"
// @Success      200  {object}  dto.PackageSearchResponse
// @Failure      500  {object}  dto.CLIErrorResponse
// @Router       /packages/search [get]
func (h *PackageHandler) Search(c *gin.Context) {
	query := c.Query("q")

	results, err := h.service.SearchPackages(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CLIErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "erreur lors de la recherche sur le registre",
		})
		return
	}

	c.JSON(http.StatusOK, dto.PackageSearchResponse{
		Count: len(results),
		Data:  results,
	})
}

// GetByNameAndVersion godoc
// @Summary      Détail d'un package et sa version
// @Description  Récupère les informations d'un package et sa version pour le CLI
// @Tags         packages
// @Produce      json
// @Param        name     path      string  true  "Nom du package"
// @Param        version  path      string  true  "Version semver ou 'latest'"
// @Success      200      {object}  dto.PackageDetailResponse
// @Failure      404      {object}  dto.CLIErrorResponse
// @Failure      500      {object}  dto.CLIErrorResponse
// @Router       /packages/{name}/{version} [get]
func (h *PackageHandler) GetByNameAndVersion(c *gin.Context) {
	name := c.Param("name")
	version := c.Param("version")

	agent, agentVer, err := h.service.GetPackageVersion(c.Request.Context(), name, version)
	if err != nil {
		if errors.Is(err, service.ErrPackageNotFound) {
			c.JSON(http.StatusNotFound, dto.CLIErrorResponse{
				Code:    "PACKAGE_NOT_FOUND",
				Message: "le package spécifié n'existe pas",
				Details: map[string]any{"name": name},
			})
			return
		}

		if errors.Is(err, service.ErrVersionNotFound) {
			c.JSON(http.StatusNotFound, dto.CLIErrorResponse{
				Code:    "VERSION_NOT_FOUND",
				Message: "la version demandée n'existe pas pour ce package",
				Details: map[string]any{"name": name, "version": version},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.CLIErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "erreur interne lors de la récupération du package",
		})
		return
	}

	// Garde-fou si le service renvoie nil sans renvoyer d'erreur explicite
	if agent == nil || agentVer == nil {
		c.JSON(http.StatusNotFound, dto.CLIErrorResponse{
			Code:    "PACKAGE_NOT_FOUND",
			Message: "package ou version introuvable",
		})
		return
	}

	c.JSON(http.StatusOK, dto.PackageDetailResponse{
		Agent:   *agent,
		Version: *agentVer,
	})
}
