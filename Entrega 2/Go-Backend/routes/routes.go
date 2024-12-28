package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter configura el router principal
func SetupRouter(router *gin.Engine) *gin.Engine {
	//router := gin.Default()

	// Inicialización de rutas
	api := router.Group("/api")
	{
		InitUserRoutes(api)
		InitCampoRoutes(api)
		InitPestTypeRoutes(api)
		InitPlagueEventRoutes(api)
		InitAlertRoutes(api)
		InitMedicionRoutes(api)
	}

	return router
}
