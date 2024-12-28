package models

// PestType representa el modelo de tipo de plaga en la base de datos
type PestType struct {
	TipoPlagaID      uint   `gorm:"primaryKey;autoIncrement" json:"tipo_plaga_id"`       // Primary key
	NombreComun      string `gorm:"type:varchar(255);not null" json:"nombre_comun"`      // Nombre común
	Descripcion      string `gorm:"type:text" json:"descripcion"`                        // Descripción
	NombreCientifico string `gorm:"type:varchar(255);not null" json:"nombre_cientifico"` // Nombre científico
}

// Esquema para la tabla en la base de datos
func (PestType) TableName() string {
	return "tipoplaga" // Nombre de la tabla en la base de datos
}
