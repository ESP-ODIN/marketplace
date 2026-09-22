package router

import (
	"marketplace/internal/handler"

	"github.com/gin-gonic/gin"
)

func registerCatalogRoutes(rg *gin.RouterGroup, h *handler.AgentHandler) {
	catalog := rg.Group("/catalog")
	{
		agents := catalog.Group("/agents")
		{
			agents.GET("", h.GetAll)
			agents.POST("", h.Create)
		}
	}
}
