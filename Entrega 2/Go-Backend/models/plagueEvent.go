package models

import (
	"time"

	"gorm.io/gorm"
)

// EventoPlaga representa el modelo de evento de plaga en la base de datos
type plagueEvent struct {
	gorm.Model
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FechaDeteccion  time.Time `json:"fecha_deteccion"`   // Fecha de detección
	Ubicacion       string    `json:"ubicacion"`         // Ubicación geográfica
	NivelSeveridad  string    `json:"nivel_severidad"`   // Nivel de severidad
	AccionesTomadas string    `json:"acciones_tomadas"`  // Acciones tomadas
	TipoPlagaID     uint      `json:"tipo_plaga_id"`     // FK de Tipo de Plaga
	CampoID         uint      `json:"campo_id"`          // FK de Campo
	RegistroVueloID uint      `json:"registro_vuelo_id"` // FK de Registro de Vuelo
	//TipoPlaga        TipoPlaga `gorm:"foreignKey:TipoPlagaID" json:"tipo_plaga"` // Relación con TipoPlaga
	//Campo            Campo     `gorm:"foreignKey:CampoID" json:"campo"`           // Relación con Campo
	//RegistroVuelo    RegistroVuelo `gorm:"foreignKey:RegistroVueloID" json:"registro_vuelo"` // Relación con RegistroVuelo
}

// Esquema
func (plagueEvent) TableName() string {
	return "evento_plaga"
}
