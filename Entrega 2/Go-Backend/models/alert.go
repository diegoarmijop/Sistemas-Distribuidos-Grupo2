package models

import (
	"time"
)

// Alert representa el modelo de alerta en la base de datos
type Alert struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"` // ID de la alerta
	Estado      string    `json:"estado"`                             // Estado de la alerta
	Descripcion string    `json:"descripcion"`                        // Descripción de la alerta
	FechaHora   time.Time `json:"fecha_hora"`                         // Fecha y hora de la alerta
	TipoAlerta  string    `json:"tipo_alerta"`                        // Tipo de alerta
	UsuarioID   uint      `json:"usuario_id"`                         // FK de Usuario
	//MedicionID    uint      `json:"medicion_id"`                         // FK de Medición
	EventoPlagaID *uint `json:"evento_plaga_id"`                     // FK de Evento de Plaga
	Usuario       *User `gorm:"foreignKey:UsuarioID" json:"usuario"` // Relación con Usuario
}

// Esquema
func (Alert) TableName() string {
	return "alert"
}
