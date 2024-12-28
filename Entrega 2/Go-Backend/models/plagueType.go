package models

import (
	"gorm.io/gorm"
)

// TipoPlaga representa el modelo de tipo de plaga en la base de datos
type PlagueType struct {
	gorm.Model
	ID               uint   `gorm:"primaryKey;autoIncrement" json:"id"` // ID principal
	NombreComun      string `json:"nombre_comun"`                       // Nombre común
	Descripcion      string `json:"descripcion"`                        // Descripción
	NombreCientifico string `json:"nombre_cientifico"`                  // Nombre científico
}

// TableName define el nombre de la tabla en la base de datos para plagueType
func (PlagueType) TableName() string {
	return "tipo_plaga"
}
