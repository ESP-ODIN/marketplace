package router

import (
	"marketplace/internal/handler"

	"github.com/gin-gonic/gin"
)

func registerHealthRoutes(r *gin.Engine, h *handler.HealthHandler) {
	health := r.Group("/health")
	{
		health.GET("/live", h.Live)
		health.GET("/ready", h.Ready)
	}
}
