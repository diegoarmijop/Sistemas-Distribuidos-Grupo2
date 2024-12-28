// main.go
package main

import (
	"github.com/gin-contrib/cors"
	"go-backend/config"
	"go-backend/routes"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Conectar a la base de datos
	config.ConnectDB()

	// Configuración del servidor
	router := gin.Default()
	gin.DefaultWriter = io.MultiWriter(os.Stdout, log.Writer())

	// Configuración de CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	// Pasar el router a SetupRouter
	routes.SetupRouter(router) // <- Aquí pasamos el router como parámetro

	// Leer el puerto desde la variable de entorno o usar el valor por defecto
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port

	log.Printf("Servidor iniciado en el puerto %s", port)

	// Inicia el servidor en el puerto especificado
	if err := router.Run(port); err != nil {
		log.Fatalf("No se pudo iniciar el servidor: %v", err)
	}
}
