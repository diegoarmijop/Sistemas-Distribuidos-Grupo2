// main.go
package main

import (
	"go-backend/config"
	"go-backend/routes"
)

func main() {
	// Se conecta la base de datos
	config.ConnectDB()

	// Configura las rutas con un router
	router := routes.SetupRouter()

	// Se elige el puerto, en este caso el 8080
	router.Run(":8080")
}
