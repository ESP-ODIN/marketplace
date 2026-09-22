package handler

import (
	"log"
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
		log.Printf("ERROR [AgentHandler.GetAll]: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch agents",
		})
		return
	}

	response := dto.ToAgentResponseList(agents)

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (h *AgentHandler) Create(c *gin.Context) {
	var req dto.CreateAgentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	createdAgent, err := h.agentService.CreateAgent(c.Request.Context(), req)
	if err != nil {
		log.Printf("ERROR [AgentHandler.Create]: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create agent",
		})
		return
	}

	response := dto.ToAgentResponse(*createdAgent)

	c.JSON(http.StatusCreated, gin.H{
		"data": response,
	})
}
