package main

import (
	"io"
	"log"
	"os"
	"sensor-dron-nodo1/config"
	"sensor-dron-nodo1/middleware"
	"sensor-dron-nodo1/routes"

	"github.com/gin-gonic/gin"
)

func main() {
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

	// Leer el puerto desde la variable de entorno o usar el valor por defecto
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = ":8081" // Valor por defecto, con el prefijo ':'
	} else if port[0] != ':' {
		port = ":" + port
	}

	log.Printf("Servidor iniciado en el puerto %s", port)

	// Inicia el servidor en el puerto especificado
	if err := router.Run(port); err != nil {
		log.Fatalf("No se pudo iniciar el servidor: %v", err)
	}
}
