package models

import (
	"time"
)

// EventoPlaga representa el modelo de evento de plaga en la base de datos
type PlagueEvent struct {
	gorm.Model
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	FechaDeteccion  time.Time  `json:"fecha_deteccion"`                          // Fecha de detección
	Ubicacion       string     `json:"ubicacion"`                                // Ubicación geográfica
	NivelSeveridad  string     `json:"nivel_severidad"`                          // Nivel de severidad
	AccionesTomadas string     `json:"acciones_tomadas"`                         // Acciones tomadas
	TipoPlagaID     uint       `json:"tipo_plaga_id"`                            // FK de Tipo de Plaga
	CampoID         uint       `json:"campo_id"`                                 // FK de Campo
	RegistroVueloID uint       `json:"registro_vuelo_id"`                        // FK de Registro de Vuelo
	TipoPlaga       PlagueType `gorm:"foreignKey:TipoPlagaID" json:"tipo_plaga"` // Relación con TipoPlaga
	Alerta          Alert      `gorm:"foreignKey:AlertaID" json:"alerta"`        // Relación con Alerta
	//Campo            Campo     `gorm:"foreignKey:CampoID" json:"campo"`           // Relación con Campo
	//RegistroVuelo    RegistroVuelo `gorm:"foreignKey:RegistroVueloID" json:"registro_vuelo"` // Relación con RegistroVuelo
}

// Esquema
func (PlagueEvent) TableName() string {
	return "evento_plaga"
}
