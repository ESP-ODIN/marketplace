package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/ESP-ODIN/authkit-go"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"

	"marketplace/internal/dto"
	"marketplace/internal/repository"
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
// @Description  Récupère tous les agents enregistrés, du plus récent au plus ancien
// @Tags         catalog
// @Produce      json
// @Success      200  {object}  dto.AgentListResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /catalog/agents [get]
func (h *AgentHandler) GetAll(c *gin.Context) {
	agents, err := h.agentService.GetAllAgents(c.Request.Context())
	if err != nil {
		log.Printf("ERROR [AgentHandler.GetAll]: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "erreur interne lors de la récupération des agents",
		})
		return
	}

	response := dto.ToAgentResponseList(agents)

	c.JSON(http.StatusOK, dto.AgentListResponse{
		Data: response,
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
		log.Printf("ERROR [AgentHandler.GetByID]: %v", err)
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
		Data: dto.ToAgentResponse(*agent),
	})
}

// Create godoc
// @Summary      Créer un agent
// @Description  Enregistre un nouvel agent sur la marketplace
// @Tags         catalog
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        agent  body      dto.CreateAgentRequest  true  "Informations de l'agent"
// @Success      201    {object}  dto.SingleAgentResponse
// @Failure      400    {object}  dto.ErrorResponse
// @Failure      401    {object}  dto.ErrorResponse
// @Failure      409    {object}  dto.ErrorResponse
// @Failure      500    {object}  dto.ErrorResponse
// @Router       /catalog/agents [post]
func (h *AgentHandler) Create(c *gin.Context) {
	identity, ok := authkit.IdentityFromContext(c.Request.Context())
	if !ok || identity.ID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "utilisateur non authentifié",
		})
		return
	}

	var req dto.CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "données invalides",
		})
		return
	}

	req.CreatorID = identity.ID

	createdAgent, err := h.agentService.CreateAgent(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrAgentAlreadyExists) {
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error: "un agent avec ce nom existe déjà",
			})
			return
		}
		if errors.Is(err, repository.ErrCreatorNotFound) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "le créateur spécifié n'existe pas",
			})
			return
		}
		log.Printf("ERROR [AgentHandler.Create]: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "erreur interne lors de la création de l'agent",
		})
		return
	}

	response := dto.ToAgentResponse(*createdAgent)

	c.JSON(http.StatusCreated, dto.SingleAgentResponse{
		Data: response,
	})
}

// GetProfile godoc
// @Summary      Agents d'un créateur
// @Description  Renvoie les agents enregistrés pour ce créateur
// @Tags         catalog
// @Produce      json
// @Param        id   path      string  true  "UUID du créateur"
// @Success      200  {object}  dto.AgentListResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /catalog/creators/{id}/agents [get]
func (h *AgentHandler) GetProfile(c *gin.Context) {
	paramID := c.Param("id")
	creatorID, err := uuid.Parse(paramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "l'identifiant fourni n'est pas un UUID valide",
		})
		return
	}

	agents, err := h.agentService.GetProfile(c.Request.Context(), creatorID)
	if err != nil {
		log.Printf("ERROR [AgentHandler.GetProfile]: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "erreur interne lors de la récupération du profil",
		})
		return
	}

	if len(agents) == 0 {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "créateur non trouvé",
		})
		return
	}

	response := dto.ToAgentResponseList(agents)

	c.JSON(http.StatusOK, dto.AgentListResponse{
		Data: response,
	})
}
