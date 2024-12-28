package services

import (
	"errors"
	"go-backend/models"

	"gorm.io/gorm"
)

type PlagueEventService struct {
	db *gorm.DB
}

func NewPlagueEventService(db *gorm.DB) *PlagueEventService {
	return &PlagueEventService{db: db}
}

func (s *PlagueEventService) Create(event *models.PlagueEvent) error {
	return s.db.Create(event).Error
}

func (s *PlagueEventService) GetByID(id uint) (*models.PlagueEvent, error) {
	var event models.PlagueEvent
	if err := s.db.Preload("TipoPlaga").Preload("Campo").First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (s *PlagueEventService) GetAll() ([]models.PlagueEvent, error) {
	var events []models.PlagueEvent
	if err := s.db.Preload("TipoPlaga").Preload("Campo").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (s *PlagueEventService) Update(event *models.PlagueEvent) error {
	if event.ID == 0 {
		return errors.New("id no válido")
	}
	return s.db.Save(event).Error
}

func (s *PlagueEventService) Delete(id uint) error {
	return s.db.Delete(&models.PlagueEvent{}, id).Error
}
