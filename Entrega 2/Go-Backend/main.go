// main.go
package main

import (
	"go-backend/config"
	"go-backend/middleware"
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
	app := gin.Default()
	gin.DefaultWriter = io.MultiWriter(os.Stdout, log.Writer())
	app.Use(middleware.CorsMiddleware()) // Registrar middleware de CORS

	// Configurar el router con las rutas definidas
	router := routes.SetupRouter()

	// Leer el puerto desde la variable de entorno o usar el valor por defecto
	port := os.Getenv("APP_PORT")

	log.Printf("Servidor iniciado en el puerto %s", port)
	if port == "" {
		port = ":8080" // Valor por defecto, con el prefijo ':'
	} else if port[0] != ':' {
		port = ":" + port
	}

	log.Printf("Servidor iniciado en el puerto %s", port)

	// Inicia el servidor en el puerto especificado
	if err := router.Run(port); err != nil {
		log.Fatalf("No se pudo iniciar el servidor: %v", err)
	}
}
