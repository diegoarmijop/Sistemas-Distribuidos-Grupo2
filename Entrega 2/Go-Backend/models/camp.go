package models

// Campo representa el modelo de campo en la base de datos
type Camp struct {
	ID          uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre      string  `json:"nombre"`                            // Nombre del campo
	Superficie  float64 `json:"superficie"`                        // Superficie del campo
	TipoCultivo string  `json:"tipo_cultivo"`                      // Tipo de cultivo
	Ubicacion   string  `json:"ubicacion"`                         // Ubicación geográfica
	SensorID    *uint   `json:"sensor_id"`                         // FK de Sensor
	Sensor      Sensor  `gorm:"foreignKey:SensorID" json:"sensor"` // Relación con Sensor
}

// Esquema
func (Camp) TableName() string {
	return "campo"
}
