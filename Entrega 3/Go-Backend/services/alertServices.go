package services

import (
	"errors"
	"go-backend/models"
	"gorm.io/gorm"
	"strings"
	"time"
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

type ResumenAlertas struct {
	Total      int    `json:"total"`
	Diferencia int    `json:"diferencia"`
	Tendencia  string `json:"tendencia"`
}

func (s *AlertService) ObtenerResumenAlertas() (*ResumenAlertas, error) {
	var totalActual int64
	var totalAnterior int64

	// Total de alertas activas
	if err := s.DB.Model(&models.Alert{}).Where("estado = ?", "Activa").Count(&totalActual).Error; err != nil {
		return nil, err
	}

	// Total de alertas de la semana anterior
	semanaAnterior := time.Now().AddDate(0, 0, -7)
	if err := s.DB.Model(&models.Alert{}).Where("fecha_hora < ? AND estado = ?", semanaAnterior, "Activa").Count(&totalAnterior).Error; err != nil {
		return nil, err
	}

	diferencia := int(totalActual - totalAnterior)
	tendencia := "down"
	if diferencia > 0 {
		tendencia = "up"
	}

	return &ResumenAlertas{
		Total:      int(totalActual),
		Diferencia: diferencia,
		Tendencia:  tendencia,
	}, nil
}

type ResolucionAlerta struct {
	SolucionAplicada string `json:"solucion_aplicada"`
	ResueltaPor      uint   `json:"resuelta_por"`
	Comentarios      string `json:"comentarios"`
	PlanAccion       string `json:"plan_accion"`
}

func (service *AlertService) ResolverAlerta(id uint, resolucion ResolucionAlerta) (*models.Alert, error) {
	var alerta models.Alert
	if err := service.DB.First(&alerta, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("alerta no encontrada")
		}
		return nil, err
	}

	ahora := time.Now()
	alerta.Estado = "Resuelta"
	alerta.SolucionAplicada = resolucion.SolucionAplicada
	alerta.FechaResolucion = &ahora
	alerta.ResueltaPor = resolucion.ResueltaPor
	alerta.Comentarios = resolucion.Comentarios
	alerta.PlanAccion = resolucion.PlanAccion

	if err := service.DB.Save(&alerta).Error; err != nil {
		return nil, err
	}

	return &alerta, nil
}

func (service *AlertService) ObtenerSugerenciasSolucion(tipoAlerta string) []string {
	// Mapa base de sugerencias por tipo individual
	sugerenciasBase := map[string][]string{
		"Humedad baja-alta": {
			"Ajustar sistema de riego",
			"Revisar drenaje del suelo",
			"Implementar sistema de ventilación",
			"Monitorear niveles de humedad cada hora",
		},
		"Temperatura alta": {
			"Activar sistema de enfriamiento",
			"Implementar sombreado temporal",
			"Aumentar ventilación",
			"Revisar aislamiento térmico",
		},
		"Temperatura baja": {
			"Activar sistema de calefacción",
			"Verificar aislamiento térmico",
			"Implementar cortinas térmicas",
		},
		"Nivel alto de insectos": {
			"Aplicar control biológico",
			"Implementar trampas para insectos",
			"Revisar barreras físicas",
			"Aplicar insecticidas orgánicos",
		},
		"Nivel de luz extremadamente alto": {
			"Instalar mallas de sombreo",
			"Ajustar orientación de cultivos",
			"Implementar sistemas de protección solar",
		},
		"Nivel de luz alto": {
			"Usar mallas de sombreo ligeras",
			"Monitorear exposición solar",
			"Ajustar horarios de exposición",
		},
	}

	// Mapa de sugerencias para combinaciones específicas
	sugerenciasCombinadas := map[string][]string{
		"Humedad baja-alta/Temperatura alta": {
			"Implementar sistema de nebulización",
			"Aumentar ventilación y humidificación",
			"Revisar y ajustar sistema de clima controlado",
			"Establecer ciclos de riego más frecuentes con menor cantidad de agua",
		},
		"Temperatura alta/Nivel alto de insectos": {
			"Implementar sistema integrado de control climático y plagas",
			"Aumentar ventilación y aplicar control biológico",
			"Revisar sellado de invernadero y aplicar control de temperatura",
			"Establecer barreras físicas con ventilación controlada",
		},
		"Humedad baja-alta/Nivel alto de insectos": {
			"Aplicar control biológico adaptado a condiciones de humedad",
			"Ajustar riego y monitorear población de insectos",
			"Implementar sistema de manejo integrado",
			"Revisar y sellar puntos de entrada de insectos",
		},
	}

	// Si existe una combinación exacta, usarla
	if sugerencias, existe := sugerenciasCombinadas[tipoAlerta]; existe {
		return sugerencias
	}

	// Si no hay combinación exacta, procesar tipos individuales
	tipos := strings.Split(tipoAlerta, "/")
	var todasSugerencias []string

	for _, tipo := range tipos {
		tipo = strings.TrimSpace(tipo)
		if sugerencias, existe := sugerenciasBase[tipo]; existe {
			todasSugerencias = append(todasSugerencias, sugerencias...)
		}
	}

	// Eliminar duplicados
	sugerenciasUnicas := make([]string, 0)
	seen := make(map[string]bool)
	for _, s := range todasSugerencias {
		if !seen[s] {
			seen[s] = true
			sugerenciasUnicas = append(sugerenciasUnicas, s)
		}
	}

	return sugerenciasUnicas
}
