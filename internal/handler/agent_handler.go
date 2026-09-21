package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"marketplace/internal/dto"
	"marketplace/internal/service"
)

type AgentHandler struct {
	agentService service.AgentService
}

func NewAgentHandler(agentService service.AgentService) *AgentHandler {
	return &AgentHandler{agentService: agentService}
}

func (h *AgentHandler) GetAll(c *gin.Context) {
	agents, err := h.agentService.GetAllAgents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les agents: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": agents,
	})
}

func (h *AgentHandler) Create(c *gin.Context) {
	var req dto.CreateAgentRequest

	// Validation du JSON avec les tags de validation du DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Données invalides",
			"details": err.Error(),
		})
		return
	}

	createdAgent, err := h.agentService.CreateAgent(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Échec de la création de l'agent : " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": createdAgent,
	})
}
