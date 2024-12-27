package routes

import (
	"go-backend/config"
	"go-backend/controllers"
	"go-backend/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	userService := services.NewUserService(config.DB)
	userController := controllers.NewUserController(userService)

	campoService := services.NewCampoService(config.DB)
	campoController := controllers.NewCampoController(campoService)

	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("/", userController.GetAllUsers)
			users.POST("/", userController.CreateUser)
			users.GET("/:id", userController.GetUser)
			users.PUT("/:id", userController.UpdateUser)
			users.DELETE("/:id", userController.DeleteUser)
		}
		// Agrupar las rutas de los campos bajo "/campos"
		campos := api.Group("/campos")
		{
			// Ruta para crear un campo
			campos.POST("/", campoController.CrearCampo)
			// Ruta para obtener todos los campos
			campos.GET("/", campoController.ObtenerTodosCampos)
			// Ruta para obtener los campos por ubicación
			campos.GET("/ubicacion/:ubicacion", campoController.ObtenerCamposPorUbicacion)
			// Ruta para actualizar el tipo de cultivo de un campo
			campos.PUT("/:id/cultivo", campoController.ActualizarCultivoCampo)
		}

	}

	return router
}
