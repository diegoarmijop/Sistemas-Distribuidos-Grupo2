package models

import "gorm.io/gorm"

// User representa el modelo de usuario en la base de datos
type User struct {
	gorm.Model
	Nombre   string `json:"nombre"`
	Email    string `json:"email" gorm:"unique"` // El correo electrónico debe ser único
	Password string `json:"password"`            // Considera encriptar la contraseña antes de guardarla
	Rol      string `json:"rol"`                 // Roles (admin, usuario, etc.)
}

// Esquema "dbo"
func (User) TableName() string {
	return "dbo.usuario" // Asumiendo que estás trabajando con el esquema dbo
}
