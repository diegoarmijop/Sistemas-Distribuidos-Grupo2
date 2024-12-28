package services

import (
	"errors"
	"go-backend/models"

	"gorm.io/gorm"
)

// PestTypeService maneja la lógica de negocio relacionada con los tipos de plaga
type PestTypeService struct {
	db *gorm.DB // Base de datos con la que interactuar
}

// NewPestTypeService crea una nueva instancia de PestTypeService
func NewPestTypeService(db *gorm.DB) *PestTypeService {
	return &PestTypeService{db: db}
}

// Create crea un nuevo tipo de plaga en la base de datos
func (s *PestTypeService) Create(pestType *models.PestType) error {
	return s.db.Create(pestType).Error
}

// GetByID obtiene un tipo de plaga desde la base de datos por su ID
func (s *PestTypeService) GetByID(id uint) (*models.PestType, error) {
	var pestType models.PestType
	// Busca el tipo de plaga en la base de datos usando el ID
	if err := s.db.First(&pestType, id).Error; err != nil {
		return nil, err
	}
	return &pestType, nil
}

// GetAll obtiene todos los tipos de plagas desde la base de datos
func (s *PestTypeService) GetAll() ([]models.PestType, error) {
	var pestTypes []models.PestType
	// Obtiene todos los tipos de plagas de la base de datos
	if err := s.db.Find(&pestTypes).Error; err != nil {
		return nil, err
	}
	return pestTypes, nil
}

// Update actualiza un tipo de plaga existente en la base de datos
func (s *PestTypeService) Update(pestType *models.PestType) error {
	// Verifica si el ID del tipo de plaga es válido
	if pestType.TipoPlagaID == 0 {
		return errors.New("id no válido")
	}
	return s.db.Save(pestType).Error
}

// Delete elimina un tipo de plaga de la base de datos por su ID
func (s *PestTypeService) Delete(id uint) error {
	return s.db.Delete(&models.PestType{}, id).Error
}
