package models

import (
	"time"

	"gorm.io/gorm"
)

// Sensor representa el modelo de sensor en la base de datos
type Sensor struct {
	gorm.Model
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TipoSensor       string    `json:"tipo_sensor"`       // Tipo de sensor
	Modelo           string    `json:"modelo"`            // Modelo del sensor
	Ubicacion        string    `json:"ubicacion"`         // Ubicación geográfica
	Estado           string    `json:"estado"`            // Estado del sensor
	FechaInstalacion time.Time `json:"fecha_instalacion"` // Fecha de instalación
}

// Esquema
func (Sensor) TableName() string {
	return "sensor"
}
