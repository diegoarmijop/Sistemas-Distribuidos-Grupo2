package routes

import (
	"sensor-dron-nodo1/config"
	"sensor-dron-nodo1/controllers"
	"sensor-dron-nodo1/services"

	"github.com/gin-gonic/gin"
)

func InitDronRoutes(api *gin.RouterGroup) {
	// Inicialización de servicios y controladores
	dronService := services.NewDronService(config.DB)
	dronController := controllers.NewDronController(dronService)

	drones := api.Group("/dron")
	{
		drones.POST("/", dronController.CrearDron)
		drones.GET("/", dronController.ObtenerTodosDrones)
		drones.PUT("/:id", dronController.ActualizarDron)
		drones.DELETE("/:id", dronController.EliminarDron)
	}
}
