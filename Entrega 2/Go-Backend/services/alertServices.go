package services

import (
	"errors"
	"go-backend/models"

	"gorm.io/gorm"
)

// AlertService gestiona las operaciones sobre las alertas
type AlertService struct {
	DB *gorm.DB
}

// NewAlertService crea una nueva instancia de AlertService
func NewAlertService(db *gorm.DB) *AlertService {
	return &AlertService{DB: db}
}

// CrearAlerta crea una nueva alerta
func (service *AlertService) CrearAlerta(alerta *models.Alert) error {
	if err := service.DB.Create(&alerta).Error; err != nil {
		return err
	}
	return nil
}

// ObtenerTodasAlertas obtiene todas las alertas
func (service *AlertService) ObtenerTodasAlertas() ([]models.Alert, error) {
	var alertas []models.Alert
	if err := service.DB.Preload("Usuario").Find(&alertas).Error; err != nil {
		return nil, err
	}
	return alertas, nil
}

// ObtenerAlertaPorID obtiene una alerta por su ID
func (service *AlertService) ObtenerAlertaPorID(id uint) (*models.Alert, error) {
	var alerta models.Alert
	if err := service.DB.Preload("Usuario").First(&alerta, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("alerta no encontrada")
		}
		return nil, err
	}
	return &alerta, nil
}

// ActualizarAlerta actualiza los detalles de una alerta
func (service *AlertService) ActualizarAlerta(id uint, updatedAlert *models.Alert) (*models.Alert, error) {
	var alerta models.Alert
	if err := service.DB.First(&alerta, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("alerta no encontrada")
		}
		return nil, err
	}

	// Actualizar los campos
	alerta.Estado = updatedAlert.Estado
	alerta.Descripcion = updatedAlert.Descripcion
	alerta.FechaHora = updatedAlert.FechaHora
	alerta.TipoAlerta = updatedAlert.TipoAlerta
	alerta.UsuarioID = updatedAlert.UsuarioID
	alerta.EventoPlagaID = updatedAlert.EventoPlagaID

	if err := service.DB.Save(&alerta).Error; err != nil {
		return nil, err
	}

	return &alerta, nil
}

// EliminarAlerta elimina una alerta por su ID
func (service *AlertService) EliminarAlerta(id uint) error {
	if err := service.DB.Delete(&models.Alert{}, id).Error; err != nil {
		return err
	}
	return nil
}

// ObtenerAlertasPorUsuarioID obtiene todas las alertas asociadas a un usuario específico
func (service *AlertService) ObtenerAlertasPorUsuarioID(usuarioID uint) ([]models.Alert, error) {
	var alertas []models.Alert
	if err := service.DB.Where("usuario_id = ?", usuarioID).Preload("Usuario").Find(&alertas).Error; err != nil {
		return nil, err
	}
	return alertas, nil
}

// ObtenerAlertasPorEventoPlagaID obtiene todas las alertas asociadas a un evento de plaga específico
func (service *AlertService) ObtenerAlertasPorEventoPlagaID(eventoPlagaID uint) ([]models.Alert, error) {
	var alertas []models.Alert
	if err := service.DB.Where("evento_plaga_id = ?", eventoPlagaID).Preload("Usuario").Find(&alertas).Error; err != nil {
		return nil, err
	}
	return alertas, nil
}
