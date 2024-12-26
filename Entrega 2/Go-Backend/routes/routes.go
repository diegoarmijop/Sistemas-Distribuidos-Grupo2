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
	}

	return router
}
