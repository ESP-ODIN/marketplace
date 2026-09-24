package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "marketplace/docs"
	"marketplace/internal/handler"
)

type Config struct {
	HealthHandler *handler.HealthHandler
	AgentHandler  *handler.AgentHandler
}

func Setup(cfg Config) *gin.Engine {
	r := gin.Default()
	_ = r.SetTrustedProxies(nil)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	registerHealthRoutes(r, cfg.HealthHandler)

	v1 := r.Group("/api/v1")
	{
		registerCatalogRoutes(v1, cfg.AgentHandler)
	}

	return r
}
