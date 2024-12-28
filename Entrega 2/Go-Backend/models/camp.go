package models

import "gorm.io/gorm"

// Campo representa el modelo de campo en la base de datos
type Camp struct {
	gorm.Model          // Esto agrega ID, CreatedAt, UpdatedAt, DeletedAt
	Nombre      string  `json:"nombre"`
	Superficie  float64 `json:"superficie"`
	TipoCultivo string  `json:"tipo_cultivo"`
	Ubicacion   string  `json:"ubicacion"`
	SensorID    *uint   `json:"sensor_id"`
	Sensor      Sensor  `gorm:"foreignKey:SensorID" json:"sensor"`
}

// Esquema
func (Camp) TableName() string {
	return "campo"
}
