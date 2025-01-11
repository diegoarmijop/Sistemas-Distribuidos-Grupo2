package models

type Dron struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Estado     string `json:"estado"`
	Modelo     string `json:"modelo"`
	Ubicacion  string `json:"ubicacion"`
	RutaID     *uint  `json:"ruta_id"` // Campo opcional
	RutaActual *Ruta  `gorm:"foreignKey:RutaID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"ruta_actual,omitempty"`
}

func (Dron) TableName() string {
	return "dron"
}
