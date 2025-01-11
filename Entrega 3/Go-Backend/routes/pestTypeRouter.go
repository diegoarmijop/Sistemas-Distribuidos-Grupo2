package routes

import (
	"go-backend/config"
	"go-backend/controllers"
	"go-backend/services"

	"github.com/gin-gonic/gin"
)

// RegisterPestTypeRoutes registra las rutas relacionadas con tipos de plagas
func InitPestTypeRoutes(api *gin.RouterGroup) {
	// Inicialización de servicios y controladores
	pestTypeService := services.NewPestTypeService(config.DB)
	pestTypeController := controllers.NewPestTypeController(pestTypeService)

	pestTypes := api.Group("/tipoPlaga")
	{
		pestTypes.POST("/", pestTypeController.Create)
		pestTypes.GET("/", pestTypeController.GetAll)
		pestTypes.GET("/:id", pestTypeController.GetByID)
		pestTypes.PUT("/:id", pestTypeController.Update)
		pestTypes.DELETE("/:id", pestTypeController.Delete)
	}
}
