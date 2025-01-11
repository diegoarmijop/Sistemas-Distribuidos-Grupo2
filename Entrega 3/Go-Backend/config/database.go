package config

import (
	"fmt"
	"go-backend/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// loadEnvVariables carga las variables de entorno desde el archivo .env
func loadEnvVariables() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error cargando archivo .env: %v", err)
	}
	return nil
}

// getDBConfig obtiene las configuraciones de la base de datos desde las variables de entorno
func getDBConfig() (string, error) {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSslmode := os.Getenv("DB_SSLMODE")
	dbTimezone := os.Getenv("DB_TIMEZONE")

	// Verificamos si alguna variable de entorno esencial está vacía
	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" || dbSslmode == "" || dbTimezone == "" {
		return "", fmt.Errorf("faltan algunas variables de entorno necesarias para la conexión")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSslmode, dbTimezone)

	return dsn, nil
}

// connectToDatabase conecta a la base de datos utilizando GORM
func connectToDatabase(dsn string) error {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error al conectar a la base de datos: %v", err)
	}
	return nil
}

// migrateDatabase realiza las migraciones automáticas en la base de datos
func migrateDatabase() error {
	// Primero migramos las tablas sin claves foráneas
	err := DB.AutoMigrate(
		&models.PestType{}, // Debe migrarse primero
	)
	if err != nil {
		return fmt.Errorf("error en la auto-migración de PestType: %v", err)
	}

	// Luego migramos las tablas que dependen de claves foráneas
	err = DB.AutoMigrate(
		&models.PlagueEvent{},
		&models.Alert{}, // Relacionado con PlagueEvent
		&models.Camp{},
		&models.Configuration{},
		&models.Sensor{},
		&models.User{},
		&models.Sensing{},
	)
	if err != nil {
		return fmt.Errorf("error en la auto-migración de PlagueEvent y otras tablas: %v", err)
	}

	return nil
}

// ConnectDB es la función principal que conecta a la base de datos y realiza la migración
func ConnectDB() {
	// Cargar variables de entorno
	if err := loadEnvVariables(); err != nil {
		log.Fatal(err)
	}

	// Obtener configuración de la base de datos
	dsn, err := getDBConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Conectar a la base de datos
	if err := connectToDatabase(dsn); err != nil {
		log.Fatal(err)
	}

	log.Println("Conexión a la base de datos exitosa")

	// Realizar migraciones
	if err := migrateDatabase(); err != nil {
		log.Fatal(err)
	}

	log.Println("Migración de base de datos completada")
}
