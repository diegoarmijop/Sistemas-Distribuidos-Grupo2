package services

import (
	"errors"
	"go-backend/models"

	"gorm.io/gorm"
)

// CampoService gestiona las operaciones sobre los campos
type CampService struct {
	DB *gorm.DB
}

// NewCampoService crea una nueva instancia de CampoService
func NewCampoService(db *gorm.DB) *CampService {
	return &CampService{DB: db}
}

// CrearCampo crea un nuevo campo si no existe en la ubicación
func (service *CampService) CrearCampo(campo *models.Camp) error {
	// Verificar si ya existe un campo en la misma ubicación
	var existingCampo models.Camp
	if err := service.DB.Where("ubicacion = ?", campo.Ubicacion).First(&existingCampo).Error; err == nil {
		return errors.New("ya existe un campo en esta ubicación")
	}

	// Guardar el nuevo campo en la base de datos
	if err := service.DB.Create(&campo).Error; err != nil {
		return err
	}
	return nil
}

// ObtenerTodosCampos obtiene todos los campos
func (service *CampService) ObtenerTodosCampos() ([]models.Camp, error) {
	var campos []models.Camp
	if err := service.DB.Find(&campos).Error; err != nil {
		return nil, err
	}
	return campos, nil
}

// ObtenerCamposPorUbicacion obtiene todos los campos de una ubicación específica
func (service *CampService) ObtenerCamposPorUbicacion(ubicacion string) ([]models.Camp, error) {
	var campos []models.Camp
	if err := service.DB.Where("ubicacion = ?", ubicacion).Find(&campos).Error; err != nil {
		return nil, err
	}
	return campos, nil
}

// ActualizarCultivoCampo actualiza el tipo de cultivo de un campo
func (service *CampService) ActualizarCultivoCampo(campoID uint, tipoCultivo string) (*models.Camp, error) {
	var campo models.Camp
	if err := service.DB.First(&campo, campoID).Error; err != nil {
		return nil, errors.New("campo no encontrado")
	}

	// Actualizar el tipo de cultivo
	campo.TipoCultivo = tipoCultivo
	if err := service.DB.Save(&campo).Error; err != nil {
		return nil, err
	}

	return &campo, nil
}
