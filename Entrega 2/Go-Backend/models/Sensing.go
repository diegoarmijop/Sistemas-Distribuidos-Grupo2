package models

import "time"

// Medicion representa el modelo de la tabla "Medición" en la base de datos
type Sensing struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"` // ID de la medición (PK)
	FechaHora   time.Time `gorm:"type:timestamp" json:"fecha_hora"`   // Fecha y hora de la medición
	Temperatura float64   `json:"temperatura"`                        // Temperatura medida
	Humedad     float64   `json:"humedad"`                            // Humedad medida
	Luminosidad float64   `json:"luminosidad"`                        // Luminosidad medida
	SensorID    uint      `json:"sensor_id"`                          // FK de Sensor
}

// Esquema
func (Sensing) TableName() string {
	return "medicion"
}
