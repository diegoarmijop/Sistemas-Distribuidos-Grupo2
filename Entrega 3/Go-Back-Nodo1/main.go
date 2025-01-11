package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sensor-dron-nodo1/config"
	"sensor-dron-nodo1/middleware"
	"sensor-dron-nodo1/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Leer el puerto desde la variable de entorno
	portEnv := os.Getenv("PORT")
	if portEnv == "" {
		portEnv = "8081" // Puerto por defecto
	}

	// Configurar logging
	gin.DefaultWriter = io.MultiWriter(os.Stdout, log.Writer())

	// Inicializar configuración primero (incluye RabbitMQ)
	config.InitConfig()

	// Verificar la conexión RabbitMQ
	if config.RabbitMQ == nil || config.RabbitMQ.Channel == nil {
		log.Fatal("Error: RabbitMQ no se inicializó correctamente")
	}

	// Conectar a la base de datos
	config.ConnectDB()

	// Configuración del servidor
	app := gin.Default()
	app.Use(middleware.CorsMiddleware()) // Registrar middleware de CORS

	// Configurar el router con las rutas definidas
	router := routes.SetupRouter()

	address := fmt.Sprintf(":%s", portEnv)
	log.Printf("Servidor iniciado en el puerto %s", address)

	// Inicia el servidor en el puerto especificado
	if err := router.Run(address); err != nil {
		log.Fatalf("No se pudo iniciar el servidor: %v", err)
	}
}
