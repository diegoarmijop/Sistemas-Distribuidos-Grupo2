package routes

import (
	"go-backend/config"
	"go-backend/controllers"
	"go-backend/services"

	"github.com/gin-gonic/gin"
)

// RegisterPlagueEventRoutes registra las rutas relacionadas con eventos de plagas
func InitPlagueEventRoutes(api *gin.RouterGroup) {
	// Inicialización de servicios y controladores
	plagueEventService := services.NewPlagueEventService(config.DB)
	plagueEventController := controllers.NewPlagueEventController(plagueEventService)

	plagueEvents := api.Group("/eventoPlagas")
	{
		plagueEvents.POST("/", plagueEventController.Create)
		plagueEvents.GET("/", plagueEventController.GetAll)
		plagueEvents.GET("/:id", plagueEventController.GetByID)
		plagueEvents.PUT("/:id", plagueEventController.Update)
		plagueEvents.DELETE("/:id", plagueEventController.Delete)
	}
}
