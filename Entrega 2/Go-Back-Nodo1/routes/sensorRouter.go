package routes

import (
	"os"
	"sensor-dron-nodo1/config"
	"sensor-dron-nodo1/controllers"
	"sensor-dron-nodo1/services"

	"github.com/gin-gonic/gin"
)

func InitSensorRoutes(api *gin.RouterGroup) {
	// Inicialización de servicios
	sensorService := services.NewSensorService(config.DB, config.RabbitMQ.Channel)
	dronService := services.NewDronService(config.DB, config.RabbitMQ.Channel)
	rutaService := services.NewRutaService(config.DB)

	// Inicialización del servicio NodoService con RabbitMQ, BaseCentral, y RutaService
	nodoService := services.NewNodoService(config.DB, config.RabbitMQ.Channel, os.Getenv("BASE_CENTRAL_URL"), rutaService)

	// Controlador de Sensor
	sensorController := controllers.NewSensorController(sensorService)

	// Iniciar el procesamiento del nodo en segundo plano
	go nodoService.ProcesarDron()

	sensores := api.Group("/sensor")
	{
		sensores.POST("/", sensorController.CrearSensor)
		sensores.GET("/", sensorController.ObtenerTodosSensores)
		sensores.PUT("/:id", sensorController.ActualizarSensor)
		sensores.DELETE("/:id", sensorController.EliminarSensor)
		sensores.POST("/publicar/:dronId", func(c *gin.Context) {
			dronID := c.Param("dronId")
			// Inicia el procesamiento del dron para este sensor
			go dronService.ProcesarDatosSensor("1", dronID) // "1" es el ID del sensor
			sensorController.PublicarDatosSensor(c)
		})
	}
}