package controllers

import (
	"net/http"
	"sensor-dron-nodo1/models"
	"sensor-dron-nodo1/services"

	"github.com/gin-gonic/gin"
)

type SensorController struct {
	SensorService *services.SensorService
}

func NewSensorController(service *services.SensorService) *SensorController {
	return &SensorController{SensorService: service}
}

// CrearSensor maneja la creación de un nuevo sensor
func (controller *SensorController) CrearSensor(c *gin.Context) {
	var sensor models.Sensor
	if err := c.ShouldBindJSON(&sensor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	if err := controller.SensorService.CrearSensor(&sensor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el sensor: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, sensor)
}

// ObtenerTodosSensores devuelve todos los sensores
func (controller *SensorController) ObtenerTodosSensores(c *gin.Context) {
	sensores, err := controller.SensorService.ObtenerTodosSensores()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener los sensores"})
		return
	}
	c.JSON(http.StatusOK, sensores)
}

// ActualizarSensor actualiza los datos de un sensor existente
func (controller *SensorController) ActualizarSensor(c *gin.Context) {
	var sensor models.Sensor
	id := c.Param("id")

	if err := controller.SensorService.ObtenerSensorPorID(id, &sensor); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sensor no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&sensor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	if err := controller.SensorService.ActualizarSensor(&sensor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el sensor: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, sensor)
}

// EliminarSensor elimina un sensor por su ID
func (controller *SensorController) EliminarSensor(c *gin.Context) {
	id := c.Param("id")

	if err := controller.SensorService.EliminarSensor(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el sensor: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sensor eliminado exitosamente"})
}

func (controller *SensorController) PublicarDatosSensor(c *gin.Context) {
	sensorID := c.Param("sensorId") // Tomar el ID del sensor de la URL

	var sensor models.Sensor
	if err := c.ShouldBindJSON(&sensor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Publicar los datos del sensor en RabbitMQ
	if err := controller.SensorService.PublicarDatos(sensor, sensorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error publicando datos del sensor: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Datos publicados exitosamente"})
}
