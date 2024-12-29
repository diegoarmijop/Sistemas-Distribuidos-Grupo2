package services

import (
	"sensor-dron-nodo1/models"

	"gorm.io/gorm"
)

type NodoService struct {
	DB *gorm.DB
}

func NewNodoService(db *gorm.DB) *NodoService {
	return &NodoService{DB: db}
}

// CrearNodo crea un nuevo nodo
func (service *NodoService) CrearNodo(nodo *models.Nodo) error {
	return service.DB.Create(nodo).Error
}

// ObtenerTodosNodos obtiene todos los nodos
func (service *NodoService) ObtenerTodosNodos() ([]models.Nodo, error) {
	var nodos []models.Nodo
	err := service.DB.Find(&nodos).Error
	return nodos, err
}

// ObtenerNodoPorID obtiene un nodo por su ID
func (service *NodoService) ObtenerNodoPorID(id string, nodo *models.Nodo) error {
	return service.DB.First(nodo, "id = ?", id).Error
}

// ActualizarNodo actualiza un nodo existente
func (service *NodoService) ActualizarNodo(nodo *models.Nodo) error {
	return service.DB.Save(nodo).Error
}

// EliminarNodo elimina un nodo por su ID
func (service *NodoService) EliminarNodo(id string) error {
	return service.DB.Delete(&models.Nodo{}, "id = ?", id).Error
}
