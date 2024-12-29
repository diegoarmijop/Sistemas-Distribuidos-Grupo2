package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	JwtKey      string
	BaseCentral string
	RabbitMQ    *RabbitMQConfig
)

func InitConfig() {
	// Cargar variables de entorno desde .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	// Inicializar RabbitMQ
	RabbitMQ = InitRabbitMQ()
	if RabbitMQ == nil || RabbitMQ.Channel == nil {
		log.Fatal("Error al inicializar RabbitMQ: conexión o canal nulo")
	}

	// Inicializar URL de la base central
	BaseCentral = os.Getenv("BASE_CENTRAL_URL")
	if BaseCentral == "" {
		log.Fatal("La URL de la base central (BASE_CENTRAL_URL) no está definida en las variables de entorno")
	}

	log.Println("Configuración inicializada correctamente")
}
