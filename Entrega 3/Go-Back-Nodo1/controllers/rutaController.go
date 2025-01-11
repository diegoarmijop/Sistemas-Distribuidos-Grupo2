package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sensor-dron-nodo1/models"
	"sensor-dron-nodo1/services"
)

type RutaController struct {
	RutaService *services.RutaService
}

func NewRutaController(service *services.RutaService) *RutaController {
	return &RutaController{RutaService: service}
}

// CrearRuta maneja la creación de una nueva ruta
func (controller *RutaController) CrearRuta(c *gin.Context) {
	var ruta models.Ruta
	if err := c.ShouldBindJSON(&ruta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	if err := controller.RutaService.CrearRuta(&ruta); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la ruta: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, ruta)
}

// ObtenerTodasRutas devuelve todas las rutas
func (controller *RutaController) ObtenerTodasRutas(c *gin.Context) {
	rutas, err := controller.RutaService.ObtenerTodasRutas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener las rutas"})
		return
	}
	c.JSON(http.StatusOK, rutas)
}

// ActualizarRuta actualiza los datos de una ruta existente
func (controller *RutaController) ActualizarRuta(c *gin.Context) {
	var ruta models.Ruta
	id := c.Param("id")

	if err := controller.RutaService.ObtenerRutaPorID(id, &ruta); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ruta no encontrada"})
		return
	}

	if err := c.ShouldBindJSON(&ruta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	if err := controller.RutaService.ActualizarRuta(&ruta); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la ruta: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, ruta)
}

// EliminarRuta elimina una ruta por su ID
func (controller *RutaController) EliminarRuta(c *gin.Context) {
	id := c.Param("id")

	if err := controller.RutaService.EliminarRuta(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la ruta: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ruta eliminada exitosamente"})
}
