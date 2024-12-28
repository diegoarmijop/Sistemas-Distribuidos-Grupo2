package routes

import (
	"go-backend/config"
	"go-backend/controllers"
	"go-backend/services"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes registra las rutas relacionadas con usuarios
func InitUserRoutes(api *gin.RouterGroup) {
	// Inicialización de servicios y controladores
	userService := services.NewUserService(config.DB)
	userController := controllers.NewUserController(userService)

	users := api.Group("/users")
	{
		users.GET("/", userController.GetAllUsers)
		users.POST("/", userController.CreateUser)
		users.GET("/:id", userController.GetUser)
		users.PUT("/:id", userController.UpdateUser)
		users.DELETE("/:id", userController.DeleteUser)
	}

	api.POST("/login", userController.Login)
}
