package routes

import (
	"sensor-dron-nodo1/config"
	"sensor-dron-nodo1/controllers"
	"sensor-dron-nodo1/services"

	"github.com/gin-gonic/gin"
)

func InitNodoRoutes(api *gin.RouterGroup) {
	// Inicialización de servicios y controladores
	nodoService := services.NewNodoService(config.DB)
	nodoController := controllers.NewNodoController(nodoService)

	nodos := api.Group("/nodo")
	{
		nodos.POST("/", nodoController.CrearNodo)
		nodos.GET("/", nodoController.ObtenerTodosNodos)
		nodos.PUT("/:id", nodoController.ActualizarNodo)
		nodos.DELETE("/:id", nodoController.EliminarNodo)
	}
}
