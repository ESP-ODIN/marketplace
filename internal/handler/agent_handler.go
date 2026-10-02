package handler

import (
	"log"
	"net/http"
	"uuid"

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

// GetAll godoc
// @Summary      Lister les agents
// @Description  Récupère la liste de tous les agents enregistrés
// @Tags         catalog
// @Produce      json
// @Success      200  {object}  dto.AgentListResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /catalog/agents [get]
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

// GetByID godoc
// @Summary      Récupérer un agent par son ID
// @Description  Renvoie les détails d'un agent spécifique à partir de son UUID
// @Tags         catalog
// @Produce      json
// @Param        id   path      string  true  "UUID de l'agent"
// @Success      200  {object}  dto.SingleAgentResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /catalog/agents/{id} [get]
func (h *AgentHandler) GetByID(c *gin.Context) {
	paramID := c.Param("id")
	agentID, err := uuid.Parse(paramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "l'identifiant fourni n'est pas un UUID valide",
		})
		return
	}

	agent, err := h.agentService.GetAgentByID(c.Request.Context(), agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "erreur interne lors de la récupération de l'agent",
		})
		return
	}

	if agent == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "agent non trouvé",
		})
		return
	}

	c.JSON(http.StatusOK, dto.SingleAgentResponse{
		Data: *agent,
	})
}

// Create godoc
// @Summary      Créer un agent
// @Description  Enregistre un nouvel agent sur la marketplace
// @Tags         catalog
// @Accept       json
// @Produce      json
// @Param        agent  body      dto.CreateAgentRequest  true  "Informations de l'agent"
// @Success      201    {object}  dto.SingleAgentResponse
// @Failure      400    {object}  dto.ErrorResponse
// @Failure      500    {object}  dto.ErrorResponse
// @Router       /catalog/agents [post]
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
