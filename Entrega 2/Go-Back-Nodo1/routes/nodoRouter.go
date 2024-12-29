package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"sensor-dron-nodo1/config"
	"sensor-dron-nodo1/controllers"
	"sensor-dron-nodo1/services"
)

func InitNodoRoutes(api *gin.RouterGroup) {
	// Inicialización de RutaService
	rutaService := services.NewRutaService(config.DB)

	// Inicialización del servicio NodoService con RabbitMQ, BaseCentral, y RutaService
	nodoService := services.NewNodoService(config.DB, config.RabbitMQ.Channel, os.Getenv("BASE_CENTRAL_URL"), rutaService)
	nodoController := controllers.NewNodoController(nodoService)

	nodos := api.Group("/nodo")
	{
		nodos.POST("/", nodoController.CrearNodo)
		nodos.GET("/", nodoController.ObtenerTodosNodos)
		nodos.PUT("/:id", nodoController.ActualizarNodo)
		nodos.DELETE("/:id", nodoController.EliminarNodo)
	}
}
