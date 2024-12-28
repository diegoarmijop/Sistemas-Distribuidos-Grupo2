package services

import (
	"sensor-dron-nodo1/models"

	"gorm.io/gorm"
)

type DronService struct {
	DB *gorm.DB
}

func NewDronService(db *gorm.DB) *DronService {
	return &DronService{DB: db}
}

// CrearDron crea un nuevo dron
func (service *DronService) CrearDron(dron *models.Dron) error {
	return service.DB.Create(&dron).Error
}

// ObtenerTodosDrones obtiene todos los drones
func (service *DronService) ObtenerTodosDrones() ([]models.Dron, error) {
	var drones []models.Dron
	err := service.DB.Find(&drones).Error
	return drones, err
}

// ObtenerDronPorID obtiene un dron por su ID
func (service *DronService) ObtenerDronPorID(id string, dron *models.Dron) error {
	return service.DB.First(dron, "id = ?", id).Error
}

// ActualizarDron actualiza un dron existente
func (service *DronService) ActualizarDron(dron *models.Dron) error {
	return service.DB.Save(dron).Error
}

// EliminarDron elimina un dron por su ID
func (service *DronService) EliminarDron(id string) error {
	return service.DB.Delete(&models.Dron{}, "id = ?", id).Error
}
