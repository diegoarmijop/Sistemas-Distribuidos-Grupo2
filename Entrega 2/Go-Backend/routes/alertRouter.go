package routes

import (
	"go-backend/config"
	"go-backend/controllers"
	"go-backend/services"

	"github.com/gin-gonic/gin"
)

// InitAlertRoutes registra las rutas relacionadas con alertas
func InitAlertRoutes(api *gin.RouterGroup) {
	// Inicialización de servicios y controladores
	alertService := services.NewAlertService(config.DB)
	alertController := controllers.NewAlertController(alertService)

	alertas := api.Group("/alertas")
	{
		alertas.POST("/", alertController.CrearAlerta)                                                // Crear una nueva alerta
		alertas.GET("/", alertController.ObtenerTodasAlertas)                                         // Obtener todas las alertas
		alertas.GET("/:id", alertController.ObtenerAlertaPorID)                                       // Obtener una alerta por ID
		alertas.PUT("/:id", alertController.ActualizarAlerta)                                         // Actualizar una alerta
		alertas.DELETE("/:id", alertController.EliminarAlerta)                                        // Eliminar una alerta
		alertas.GET("/usuario/:usuario_id", alertController.ObtenerAlertasPorUsuarioID)               // Obtener alertas por ID de usuario
		alertas.GET("/evento_plaga/:evento_plaga_id", alertController.ObtenerAlertasPorEventoPlagaID) // Obtener alertas por ID de evento de plaga
	}
}
