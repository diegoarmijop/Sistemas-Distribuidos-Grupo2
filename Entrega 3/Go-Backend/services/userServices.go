package services

import (
	"go-backend/models"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

// services/user_service.go

func (s *UserService) LoginUser(email, password string) (*models.User, error) {
	var user models.User

	// Buscar usuario por email
	if err := s.DB.Where("email = ? AND password = ?", email, password).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.DB.Create(user).Error
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := s.DB.Find(&users)
	return users, result.Error
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) DeleteUser(id uint) error {
	// Primero verificamos que el usuario existe
	var user models.User
	if err := s.DB.First(&user, id).Error; err != nil {
		return err
	}

	// Realizar el borrado suave (soft delete) ya que estamos usando gorm.Model
	return s.DB.Delete(&user).Error
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(id uint, userData *models.User) error {
	// Primero verificamos que el usuario existe
	var existingUser models.User
	if err := s.DB.First(&existingUser, id).Error; err != nil {
		return err
	}

	// Actualizamos solo nombre y email
	updates := map[string]interface{}{
		"nombre": userData.Nombre,
		"email":  userData.Email,
	}

	// Realizar la actualización
	return s.DB.Model(&existingUser).Updates(updates).Error
}
