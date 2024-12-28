package services

import (
	"errors"
	"go-backend/models"

	"gorm.io/gorm"
)

// MedicionService gestiona las operaciones sobre las mediciones
type MedicionService struct {
	DB *gorm.DB
}

// NewMedicionService crea una nueva instancia de MedicionService
func NewMedicionService(db *gorm.DB) *MedicionService {
	return &MedicionService{DB: db}
}

// CrearMedicion crea una nueva medición en la base de datos
func (service *MedicionService) CrearMedicion(medicion *models.Sensing) error {
	// Guardar la nueva medición en la base de datos
	if err := service.DB.Create(&medicion).Error; err != nil {
		return err
	}
	return nil
}

// ObtenerTodasMediciones obtiene todas las mediciones
func (service *MedicionService) ObtenerTodasMediciones() ([]models.Sensing, error) {
	var mediciones []models.Sensing
	if err := service.DB.Find(&mediciones).Error; err != nil {
		return nil, err
	}
	return mediciones, nil
}

// ObtenerMedicionesPorSensor obtiene todas las mediciones de un sensor específico
func (service *MedicionService) ObtenerMedicionesPorSensor(sensorID uint) ([]models.Sensing, error) {
	var mediciones []models.Sensing
	if err := service.DB.Where("sensor_id = ?", sensorID).Find(&mediciones).Error; err != nil {
		return nil, err
	}
	return mediciones, nil
}

// ActualizarMedicion actualiza una medición específica
func (service *MedicionService) ActualizarMedicion(medicionID uint, nuevaMedicion *models.Sensing) (*models.Sensing, error) {
	var medicion models.Sensing
	if err := service.DB.First(&medicion, medicionID).Error; err != nil {
		return nil, errors.New("medición no encontrada")
	}

	// Actualizar los campos de la medición
	medicion.FechaHora = nuevaMedicion.FechaHora
	medicion.Temperatura = nuevaMedicion.Temperatura
	medicion.Humedad = nuevaMedicion.Humedad
	medicion.Luminosidad = nuevaMedicion.Luminosidad
	medicion.SensorID = nuevaMedicion.SensorID

	if err := service.DB.Save(&medicion).Error; err != nil {
		return nil, err
	}

	return &medicion, nil
}

// EliminarMedicion elimina una medición por su ID
func (service *MedicionService) EliminarMedicion(medicionID uint) error {
	if err := service.DB.Delete(&models.Sensing{}, medicionID).Error; err != nil {
		return err
	}
	return nil
}
