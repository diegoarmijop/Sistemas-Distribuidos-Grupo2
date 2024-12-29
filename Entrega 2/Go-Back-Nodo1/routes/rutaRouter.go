package routes

import (
	"github.com/gin-gonic/gin"
	"sensor-dron-nodo1/config"
	"sensor-dron-nodo1/controllers"
	"sensor-dron-nodo1/services"
)

func InitRutaRoutes(api *gin.RouterGroup) {
	// Inicialización de servicios y controladores
	rutaService := services.NewRutaService(config.DB)
	rutaController := controllers.NewRutaController(rutaService)

	rutas := api.Group("/ruta")
	{
		rutas.POST("/", rutaController.CrearRuta)
		rutas.GET("/", rutaController.ObtenerTodasRutas)
		rutas.PUT("/:id", rutaController.ActualizarRuta)
		rutas.DELETE("/:id", rutaController.EliminarRuta)
	}
}
