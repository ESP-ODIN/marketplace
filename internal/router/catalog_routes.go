package router

import (
	"github.com/ESP-ODIN/authkit-go"
	"github.com/gin-gonic/gin"

	"marketplace/internal/handler"
	"marketplace/internal/middleware"
)

func registerCatalogRoutes(rg *gin.RouterGroup, h *handler.AgentHandler, authenticator *authkit.Authenticator) {
	catalog := rg.Group("/catalog")
	{
		agents := catalog.Group("/agents")
		{
			agents.GET("", h.GetAll)
			agents.GET("/:id", h.GetByID)
			agents.POST("", middleware.RequireAuth(authenticator), h.Create)
		}
		creators := catalog.Group("/creators")
		{
			creators.GET("/:id/agents", h.GetProfile)
		}
	}
}
