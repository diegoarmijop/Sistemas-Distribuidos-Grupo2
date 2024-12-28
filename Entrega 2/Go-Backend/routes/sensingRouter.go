package routes

import (
	"go-backend/config"
	"go-backend/controllers"
	"go-backend/services"

	"github.com/gin-gonic/gin"
)

// InitMedicionRoutes registra las rutas relacionadas con mediciones
func InitMedicionRoutes(api *gin.RouterGroup) {
	// Inicialización de servicios y controladores
	medicionService := services.NewMedicionService(config.DB)
	medicionController := controllers.NewMedicionController(medicionService)

	mediciones := api.Group("/mediciones")
	{
		mediciones.POST("/", medicionController.CrearMedicion)
		mediciones.GET("/", medicionController.ObtenerTodasMediciones)
		mediciones.GET("/sensor/:sensor_id", medicionController.ObtenerMedicionesPorSensor)
		mediciones.PUT("/:id", medicionController.ActualizarMedicion)
		mediciones.DELETE("/:id", medicionController.EliminarMedicion)
	}
}
