package services

import (
	"errors"
	"go-backend/models"
	"math"
	"time"

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

type ResumenControl struct {
	Efectividad float64 `json:"efectividad"`
	Total       int     `json:"total"`
	Tendencia   string  `json:"tendencia"`
}

func (s *PlagueEventService) ObtenerResumenControl() (*ResumenControl, error) {
	var totalEventos int64
	var eventosControlados int64

	if err := s.db.Model(&models.PlagueEvent{}).Count(&totalEventos).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.PlagueEvent{}).Where("nivel_severidad = ?", "Controlado").Count(&eventosControlados).Error; err != nil {
		return nil, err
	}

	efectividad := float64(0)
	if totalEventos > 0 {
		efectividad = (float64(eventosControlados) / float64(totalEventos)) * 100
	}

	return &ResumenControl{
		Efectividad: math.Round(efectividad*100) / 100,
		Total:       int(totalEventos),
		Tendencia:   "up",
	}, nil
}

type EstadisticasEventos struct {
	Labels []string `json:"labels"` // Fechas/meses
	Datos  []int    `json:"datos"`  // Cantidad de eventos
	Total  int      `json:"total"`
}

func (s *PlagueEventService) ObtenerEstadisticas() (*EstadisticasEventos, error) {
	var resultados []struct {
		Fecha  time.Time `json:"fecha"`
		Conteo int       `json:"conteo"`
	}

	// Obtener datos de los últimos 6 meses
	seisMesesAtras := time.Now().AddDate(0, -6, 0)

	err := s.db.Model(&models.PlagueEvent{}).
		Select("DATE(fecha_deteccion) as fecha, count(*) as conteo").
		Where("fecha_deteccion >= ?", seisMesesAtras).
		Group("DATE(fecha_deteccion)").
		Order("fecha").
		Find(&resultados).Error

	if err != nil {
		return nil, err
	}

	// Procesar resultados
	labels := make([]string, len(resultados))
	datos := make([]int, len(resultados))
	total := 0

	for i, r := range resultados {
		labels[i] = r.Fecha.Format("02/01")
		datos[i] = r.Conteo
		total += r.Conteo
	}

	return &EstadisticasEventos{
		Labels: labels,
		Datos:  datos,
		Total:  total,
	}, nil
}
