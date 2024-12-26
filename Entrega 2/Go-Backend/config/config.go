package config

import (
	"log"

	"github.com/joho/godotenv"
)

var JwtKey string

func InitConfig() {
	// Cargar variables de entorno desde .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	// Asignar la clave secreta de JWT desde las variables de entorno
	//JwtKey = os.Getenv("JWT_SECRET")
	//if JwtKey == "" {
	//    log.Fatal("La clave secreta JWT no está definida en las variables de entorno")
	//}

	// Depuración: Imprimir la clave JWT para confirmar que se carga correctamente
	//log.Println("JWT Key cargada:", JwtKey)
}
