package router

import (
	"marketplace/internal/handler"

	"github.com/gin-gonic/gin"
)

type Config struct {
	HealthHandler *handler.HealthHandler
	AgentHandler  *handler.AgentHandler
}

func Setup(cfg Config) *gin.Engine {
	r := gin.Default()
	_ = r.SetTrustedProxies(nil)

	r.GET("/health", cfg.HealthHandler.Check)

	v1 := r.Group("/api/v1")
	{
		agents := v1.Group("/agents")
		{
			agents.GET("", cfg.AgentHandler.GetAll)
			agents.POST("", cfg.AgentHandler.Create)
		}
	}

	return r
}
