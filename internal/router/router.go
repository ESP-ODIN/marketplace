package router

import (
	"marketplace/internal/handler"

	"github.com/gin-gonic/gin"
)

type Config struct {
	HealthHandler *handler.HealthHandler
	AgentHandler  *handler.AgentHandler
	// Plus tard : PackageHandler *handler.PackageHandler
}

func Setup(cfg Config) *gin.Engine {
	r := gin.Default()
	_ = r.SetTrustedProxies(nil)

	// Route racine /health
	registerHealthRoutes(r, cfg.HealthHandler)

	// Préfixe v1
	v1 := r.Group("/api/v1")
	{
		registerCatalogRoutes(v1, cfg.AgentHandler)
		// Plus tard : registerPackageRoutes(v1, cfg.PackageHandler)
	}

	return r
}
