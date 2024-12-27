package models

import "gorm.io/gorm"

// Configuracion representa el modelo de configuración en la base de datos
type Configuration struct {
	gorm.Model
	ID               uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	UmbraTemp        float64 `json:"umbral_temp"`                         // Umbral de temperatura
	UmbraHumedad     float64 `json:"umbral_humedad"`                      // Umbral de humedad
	UmbraLuminosidad float64 `json:"umbral_luminocidad"`                  // Umbral de luminosidad
	UsuarioID        uint    `json:"usuario_id"`                          // FK de Usuario
	Usuario          User    `gorm:"foreignKey:UsuarioID" json:"usuario"` // Relación con Usuario
}

// Esquema
func (Configuration) TableName() string {
	return "configuracion"
}
