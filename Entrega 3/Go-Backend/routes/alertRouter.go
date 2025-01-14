package routes

import (
	"github.com/gin-gonic/gin"
	"go-backend/config"
	"go-backend/controllers"
	"go-backend/services"
)

func InitAlertRoutes(api *gin.RouterGroup) {
	alertService := services.NewAlertService(config.DB)
	alertController := controllers.NewAlertController(alertService)

	alertas := api.Group("/alertas")
	{
		// Rutas existentes
		alertas.POST("/", alertController.CrearAlerta)
		alertas.GET("/", alertController.ObtenerTodasAlertas)
		alertas.GET("/:id", alertController.ObtenerAlertaPorID)
		alertas.PUT("/:id", alertController.ActualizarAlerta)
		alertas.DELETE("/:id", alertController.EliminarAlerta)
		alertas.GET("/usuario/:usuario_id", alertController.ObtenerAlertasPorUsuarioID)
		alertas.GET("/evento_plaga/:evento_plaga_id", alertController.ObtenerAlertasPorEventoPlagaID)
		alertas.GET("/resumen", alertController.ObtenerResumenAlertas)

		// Nuevas rutas
		alertas.POST("/:id/resolver", alertController.ResolverAlerta)
		alertas.GET("/sugerencias", alertController.ObtenerSugerenciasSolucion)
	}
}
