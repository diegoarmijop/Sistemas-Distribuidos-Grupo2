package routes

import (
	"sensor-dron-nodo1/config"
	"sensor-dron-nodo1/controllers"
	"sensor-dron-nodo1/services"

	"github.com/gin-gonic/gin"
)

func InitSensorRoutes(api *gin.RouterGroup) {
	// Inicialización de servicios y controladores
	sensorService := services.NewSensorService(config.DB)
	sensorController := controllers.NewSensorController(sensorService)

	sensores := api.Group("/sensor")
	{
		sensores.POST("/", sensorController.CrearSensor)
		sensores.GET("/", sensorController.ObtenerTodosSensores)
		sensores.PUT("/:id", sensorController.ActualizarSensor)
		sensores.DELETE("/:id", sensorController.EliminarSensor)
	}
}
