package routes

import (
	"go-backend/config"
	"go-backend/controllers"
	"go-backend/services"

	"github.com/gin-gonic/gin"
)

// RegisterCampoRoutes registra las rutas relacionadas con campos
func InitCampoRoutes(api *gin.RouterGroup) {
	// Inicialización de servicios y controladores
	campoService := services.NewCampoService(config.DB)
	campoController := controllers.NewCampoController(campoService)

	campos := api.Group("/campos")
	{
		campos.POST("/", campoController.CrearCampo)
		campos.GET("/", campoController.ObtenerTodosCampos)
		campos.GET("/ubicacion/:ubicacion", campoController.ObtenerCamposPorUbicacion)
		campos.PUT("/:id/cultivo", campoController.ActualizarCultivoCampo)
	}
}
