package services

import (
	"sensor-dron-nodo1/models"

	"gorm.io/gorm"
)

type SensorService struct {
	DB *gorm.DB
}

func NewSensorService(db *gorm.DB) *SensorService {
	return &SensorService{DB: db}
}

// CrearSensor crea un nuevo sensor
func (service *SensorService) CrearSensor(sensor *models.Sensor) error {
	return service.DB.Create(sensor).Error
}

// ObtenerTodosSensores obtiene todos los sensores
func (service *SensorService) ObtenerTodosSensores() ([]models.Sensor, error) {
	var sensores []models.Sensor
	err := service.DB.Find(&sensores).Error
	return sensores, err
}

// ObtenerSensorPorID obtiene un sensor por su ID
func (service *SensorService) ObtenerSensorPorID(id string, sensor *models.Sensor) error {
	return service.DB.First(sensor, "id = ?", id).Error
}

// ActualizarSensor actualiza un sensor existente
func (service *SensorService) ActualizarSensor(sensor *models.Sensor) error {
	return service.DB.Save(sensor).Error
}

// EliminarSensor elimina un sensor por su ID
func (service *SensorService) EliminarSensor(id string) error {
	return service.DB.Delete(&models.Sensor{}, "id = ?", id).Error
}
