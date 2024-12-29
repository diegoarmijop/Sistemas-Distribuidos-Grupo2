package services

import (
	"sensor-dron-nodo1/models"

	"gorm.io/gorm"
)

type RutaService struct {
	DB *gorm.DB
}

func NewRutaService(db *gorm.DB) *RutaService {
	return &RutaService{DB: db}
}

// CrearRuta crea una nueva ruta
func (service *RutaService) CrearRuta(ruta *models.Ruta) error {
	return service.DB.Create(ruta).Error
}

// ObtenerTodasRutas obtiene todas las rutas
func (service *RutaService) ObtenerTodasRutas() ([]models.Ruta, error) {
	var rutas []models.Ruta
	err := service.DB.Find(&rutas).Error
	return rutas, err
}

// ObtenerRutaPorID obtiene una ruta por su ID
func (service *RutaService) ObtenerRutaPorID(id string, ruta *models.Ruta) error {
	return service.DB.First(ruta, "id = ?", id).Error
}

// ActualizarRuta actualiza una ruta existente
func (service *RutaService) ActualizarRuta(ruta *models.Ruta) error {
	return service.DB.Save(ruta).Error
}

// EliminarRuta elimina una ruta por su ID
func (service *RutaService) EliminarRuta(id string) error {
	return service.DB.Delete(&models.Ruta{}, "id = ?", id).Error
}
