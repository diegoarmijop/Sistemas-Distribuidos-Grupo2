package routes

import (
	"go-backend/config"
	"go-backend/controllers"
	"go-backend/services"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Configuración de CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Inicialización de servicios
	userService := services.NewUserService(config.DB)
	campoService := services.NewCampoService(config.DB)
	pestTypeService := services.NewPestTypeService(config.DB)
	plagueEventService := services.NewPlagueEventService(config.DB)

	// Inicialización de controladores
	userController := controllers.NewUserController(userService)
	campoController := controllers.NewCampoController(campoService)
	pestTypeController := controllers.NewPestTypeController(pestTypeService)
	plagueEventController := controllers.NewPlagueEventController(plagueEventService)

	api := router.Group("/api")
	{
		// Rutas de autenticación
		api.POST("/login", userController.Login)

		// Rutas de usuarios
		users := api.Group("/users")
		{
			users.GET("/", userController.GetAllUsers)
			users.POST("/", userController.CreateUser)
			users.GET("/:id", userController.GetUser)
			users.PUT("/:id", userController.UpdateUser)
			users.DELETE("/:id", userController.DeleteUser)
		}

		// Rutas de campos
		campos := api.Group("/campos")
		{
			campos.POST("/", campoController.CrearCampo)
			campos.GET("/", campoController.ObtenerTodosCampos)
			campos.GET("/ubicacion/:ubicacion", campoController.ObtenerCamposPorUbicacion)
			campos.PUT("/:id/cultivo", campoController.ActualizarCultivoCampo)
		}

		// Rutas para TipoPlaga
		pestTypes := api.Group("/tipoPlaga")
		{
			pestTypes.POST("/", pestTypeController.Create)
			pestTypes.GET("/", pestTypeController.GetAll)
			pestTypes.GET("/:id", pestTypeController.GetByID)
			pestTypes.PUT("/:id", pestTypeController.Update)
			pestTypes.DELETE("/:id", pestTypeController.Delete)
		}

		// Rutas para EventoPlaga
		plagueEvents := api.Group("/eventoPlagas")
		{
			plagueEvents.POST("/", plagueEventController.Create)
			plagueEvents.GET("/", plagueEventController.GetAll)
			plagueEvents.GET("/:id", plagueEventController.GetByID)
			plagueEvents.PUT("/:id", plagueEventController.Update)
			plagueEvents.DELETE("/:id", plagueEventController.Delete)
		}
	}

	return router
}
