package services

import (
	"errors"
	"go-backend/models"
	"gorm.io/gorm"
	"time"
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

// services/camp_service.go
type ResumenCampos struct {
	Total     int    `json:"total"`
	Nuevos    int    `json:"nuevos"`
	Tendencia string `json:"tendencia"`
}

// Agregar el método para obtener el resumen
func (service *CampService) ObtenerResumenCampos() (*ResumenCampos, error) {
	var total int64
	var nuevos int64

	// Contar total de campos
	if err := service.DB.Model(&models.Camp{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// Contar campos nuevos este mes
	primerDiaMes := time.Now().AddDate(0, 0, -30)
	if err := service.DB.Model(&models.Camp{}).Where("created_at >= ?", primerDiaMes).Count(&nuevos).Error; err != nil {
		return nil, err
	}

	tendencia := "down"
	if nuevos > 0 {
		tendencia = "up"
	}

	return &ResumenCampos{
		Total:     int(total),
		Nuevos:    int(nuevos),
		Tendencia: tendencia,
	}, nil
}
