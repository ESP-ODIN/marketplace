// Dans internal/router/package_routes.go
package router

import (
	"marketplace/internal/handler"

	"github.com/gin-gonic/gin"
)

func registerPackageRoutes(rg *gin.RouterGroup, h *handler.PackageHandler) {
	packages := rg.Group("/packages")
	{
		packages.GET("/search", h.Search)
		packages.GET("/:name/:version", h.GetByNameAndVersion)
	}
}
