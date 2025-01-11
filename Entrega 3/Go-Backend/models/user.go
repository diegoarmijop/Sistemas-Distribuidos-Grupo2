package models

import "gorm.io/gorm"

// User representa el modelo de usuario en la base de datos
type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre   string `json:"nombre"`
	Email    string `json:"email" gorm:"unique"` // El correo electrónico debe ser único
	Password string `json:"password"`            // Considera encriptar la contraseña antes de guardarla
	Rol      string `json:"rol"`                 // Roles (admin, usuario, etc.)
}

// Esquema
func (User) TableName() string {
	return "usuario"
}
