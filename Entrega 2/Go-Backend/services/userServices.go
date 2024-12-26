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

func (s *UserService) UpdateUser(user *models.User) error {
	return s.DB.Save(user).Error
}

func (s *UserService) DeleteUser(id uint) error {
	return s.DB.Delete(&models.User{}, id).Error
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
