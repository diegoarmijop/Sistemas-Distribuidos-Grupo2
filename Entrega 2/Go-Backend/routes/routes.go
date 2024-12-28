package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter configura el router principal
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Inicialización de rutas
	api := router.Group("/api")
	{
		InitUserRoutes(api)
		InitCampoRoutes(api)
		InitPestTypeRoutes(api)
		InitPlagueEventRoutes(api)
	}

	return router
}
