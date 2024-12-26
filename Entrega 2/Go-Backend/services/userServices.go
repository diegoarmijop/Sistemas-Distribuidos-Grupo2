package services

import (
	"go-backend/models"

	"gorm.io/gorm"
)

// userServices estructura para manejar los servicios de usuarios
type UserServices struct {
	DB *gorm.DB
}

// NewUsuarioService crea una nueva instancia de userServices
func NewUsuarioService(db *gorm.DB) *UserServices {
	return &UserServices{DB: db}
}

// ObtenerUsuarioPorID obtiene un usuario por su ID
func (us *UserServices) ObtenerUsuarioPorID(userID int) (*models.User, error) {
	var usuario models.User
	// Consultar la base de datos buscando un usuario por ID
	if err := us.DB.Where("id = ?", userID).First(&usuario).Error; err != nil {
		return nil, err
	}
	return &usuario, nil
}

// CrearUsuario crea un nuevo usuario en la base de datos
func (us *UserServices) CrearUsuario(usuario models.User) (*models.User, error) {
	// Crear un usuario en la base de datos
	if err := us.DB.Create(&usuario).Error; err != nil {
		return nil, err
	}
	return &usuario, nil
}

// ObtenerUsuarios obtiene todos los usuarios de la base de datos
func (us *UserServices) ObtenerUsuarios() ([]models.User, error) {
	var usuarios []models.User
	// Consultar todos los usuarios
	if err := us.DB.Find(&usuarios).Error; err != nil {
		return nil, err
	}
	return usuarios, nil
}

// ActualizarUsuario actualiza los datos de un usuario en la base de datos
func (us *UserServices) ActualizarUsuario(userID int, usuario models.User) (*models.User, error) {
	var existingUser models.User
	// Buscar el usuario existente por ID
	if err := us.DB.Where("id = ?", userID).First(&existingUser).Error; err != nil {
		return nil, err
	}
	// Actualizar los campos necesarios
	existingUser.Nombre = usuario.Nombre
	existingUser.Email = usuario.Email
	existingUser.Password = usuario.Password
	existingUser.Rol = usuario.Rol

	// Guardar los cambios
	if err := us.DB.Save(&existingUser).Error; err != nil {
		return nil, err
	}

	return &existingUser, nil
}

// EliminarUsuario elimina un usuario por su ID de la base de datos
func (us *UserServices) EliminarUsuario(userID int) error {
	var usuario models.User
	// Buscar el usuario por ID
	if err := us.DB.Where("id = ?", userID).First(&usuario).Error; err != nil {
		return err
	}
	// Eliminar el usuario
	if err := us.DB.Delete(&usuario).Error; err != nil {
		return err
	}
	return nil
}
