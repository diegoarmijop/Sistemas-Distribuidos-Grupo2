package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	config := cors.Config{
		AllowAllOrigins:  true, // Permitir solicitudes desde cualquier origen
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Origin"},
		ExposeHeaders:    []string{"Pagination-Count"}, // Si necesitas exponer algún encabezado
		AllowCredentials: false,                        // Ajustar según sea necesario
	}

	return cors.New(config)
}
