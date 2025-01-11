package controllers

import (
	"go-backend/models"
	"go-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AlertController maneja las peticiones relacionadas con las alertas
type AlertController struct {
	AlertService *services.AlertService
}

// NewAlertController crea una nueva instancia de AlertController
func NewAlertController(service *services.AlertService) *AlertController {
	return &AlertController{AlertService: service}
}

// CrearAlerta maneja la creación de una nueva alerta
func (controller *AlertController) CrearAlerta(c *gin.Context) {
	var alerta models.Alert
	if err := c.ShouldBindJSON(&alerta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := controller.AlertService.CrearAlerta(&alerta); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alerta)
}

// ObtenerTodasAlertas maneja la obtención de todas las alertas
func (controller *AlertController) ObtenerTodasAlertas(c *gin.Context) {
	alertas, err := controller.AlertService.ObtenerTodasAlertas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las alertas"})
		return
	}

	c.JSON(http.StatusOK, alertas)
}

// ObtenerAlertaPorID maneja la obtención de una alerta por su ID
func (controller *AlertController) ObtenerAlertaPorID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	alerta, err := controller.AlertService.ObtenerAlertaPorID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alerta)
}

// ActualizarAlerta maneja la actualización de una alerta
func (controller *AlertController) ActualizarAlerta(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var updatedAlert models.Alert
	if err := c.ShouldBindJSON(&updatedAlert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	alerta, err := controller.AlertService.ActualizarAlerta(uint(id), &updatedAlert)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alerta)
}

// EliminarAlerta maneja la eliminación de una alerta
func (controller *AlertController) EliminarAlerta(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := controller.AlertService.EliminarAlerta(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alerta eliminada exitosamente"})
}

// ObtenerAlertasPorUsuarioID maneja la obtención de alertas por ID de usuario
func (controller *AlertController) ObtenerAlertasPorUsuarioID(c *gin.Context) {
	usuarioIDStr := c.Param("usuario_id")
	usuarioID, err := strconv.ParseUint(usuarioIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	alertas, err := controller.AlertService.ObtenerAlertasPorUsuarioID(uint(usuarioID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alertas)
}

// ObtenerAlertasPorEventoPlagaID maneja la obtención de alertas por ID de evento de plaga
func (controller *AlertController) ObtenerAlertasPorEventoPlagaID(c *gin.Context) {
	eventoPlagaIDStr := c.Param("evento_plaga_id")
	eventoPlagaID, err := strconv.ParseUint(eventoPlagaIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de evento de plaga inválido"})
		return
	}

	alertas, err := controller.AlertService.ObtenerAlertasPorEventoPlagaID(uint(eventoPlagaID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alertas)
}

func (c *AlertController) ObtenerResumenAlertas(ctx *gin.Context) {
	resumen, err := c.AlertService.ObtenerResumenAlertas()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener resumen de alertas",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, resumen)
}
