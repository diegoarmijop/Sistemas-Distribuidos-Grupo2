package controllers

import (
	"go-backend/models"
	"go-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MedicionController maneja las peticiones relacionadas con las mediciones
type MedicionController struct {
	MedicionService *services.MedicionService
}

// NewMedicionController crea una nueva instancia de MedicionController
func NewMedicionController(service *services.MedicionService) *MedicionController {
	return &MedicionController{MedicionService: service}
}

// CrearMedicion maneja la creación de una nueva medición
func (controller *MedicionController) CrearMedicion(c *gin.Context) {
	var medicion models.Sensing
	if err := c.ShouldBindJSON(&medicion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Usar el servicio para crear la medición
	if err := controller.MedicionService.CrearMedicion(&medicion); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, medicion)
}

// ObtenerTodasMediciones maneja la obtención de todas las mediciones
func (controller *MedicionController) ObtenerTodasMediciones(c *gin.Context) {
	mediciones, err := controller.MedicionService.ObtenerTodasMediciones()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las mediciones"})
		return
	}

	c.JSON(http.StatusOK, mediciones)
}

// ObtenerMedicionesPorSensor obtiene mediciones por sensor específico
func (controller *MedicionController) ObtenerMedicionesPorSensor(c *gin.Context) {
	// Obtener el sensor_id de la URL y convertirlo a uint
	sensorIDStr := c.Param("sensor_id")
	sensorID, err := strconv.ParseUint(sensorIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de sensor inválido"})
		return
	}

	// Usar el servicio para obtener las mediciones
	mediciones, err := controller.MedicionService.ObtenerMedicionesPorSensor(uint(sensorID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mediciones)
}

// ActualizarMedicion permite actualizar una medición
func (controller *MedicionController) ActualizarMedicion(c *gin.Context) {
	// Obtener el medicion_id de la URL y convertirlo a uint
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Crear un struct para recibir los datos actualizados
	var nuevaMedicion models.Sensing
	if err := c.ShouldBindJSON(&nuevaMedicion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Usar el servicio para actualizar la medición
	medicion, err := controller.MedicionService.ActualizarMedicion(uint(id), &nuevaMedicion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, medicion)
}

// EliminarMedicion maneja la eliminación de una medición
func (controller *MedicionController) EliminarMedicion(c *gin.Context) {
	// Obtener el medicion_id de la URL y convertirlo a uint
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Usar el servicio para eliminar la medición
	if err := controller.MedicionService.EliminarMedicion(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medición eliminada correctamente"})
}
