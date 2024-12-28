package models

type Nodo struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Estado    string `json:"estado"`
	Ubicacion string `json:"ubicacion"`
}

// Esquema
func (Nodo) TableName() string {
	return "nodo"
}
