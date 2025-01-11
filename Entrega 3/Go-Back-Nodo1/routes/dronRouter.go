package routes

import (
	"log"
	"sensor-dron-nodo1/config"
	"sensor-dron-nodo1/controllers"
	"sensor-dron-nodo1/services"

	"github.com/gin-gonic/gin"
)

func InitDronRoutes(api *gin.RouterGroup) {
	if config.RabbitMQ == nil || config.RabbitMQ.Channel == nil {
		log.Fatal("RabbitMQ no está inicializado correctamente")
	}

	dronService := services.NewDronService(config.DB, config.RabbitMQ.Channel)
	dronController := controllers.NewDronController(dronService)

	drones := api.Group("/dron")
	{
		drones.POST("/", dronController.CrearDron)
		drones.GET("/", dronController.ObtenerTodosDrones)
		drones.PUT("/:id", dronController.ActualizarDron)
		drones.DELETE("/:id", dronController.EliminarDron)
	}
}
