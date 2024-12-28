package models

type Dron struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Estado     string `json:"estado"`
	Modelo     string `json:"modelo"`
	Ubicacion  string `json:"ubicacion"`
	RutaID     uint   `json:"ruta_id"` // Clave foránea explícita
	RutaActual Ruta   `gorm:"foreignKey:RutaID" json:"ruta_actual"`
}

// Esquema
func (Dron) TableName() string {
	return "dron"
}
