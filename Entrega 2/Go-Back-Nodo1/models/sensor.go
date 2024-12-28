package models

type Sensor struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Temperatura string `json:"temperatura"`
	Humedad     string `json:"humedad"`
	Insectos    string `json:"insectos"`
	Luz         string `json:"luz"`
}

// Esquema
func (Sensor) TableName() string {
	return "sensor"
}
