package models

import (
	"time"
)

type Ruta struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FechaHoraInicio  time.Time `json:"fecha_hora_inicio"`
	FechaHoraTermino time.Time `json:"fecha_hora_termino"`
}

// Esquema
func (Ruta) TableName() string {
	return "ruta"
}
