package controllers

import (
	"go-backend/models"
	"go-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PlagueTypeController maneja las peticiones relacionadas con los tipos de plaga
type PlagueTypeController struct {
	PlagueTypeService *services.PlagueTypeService
}

// NewPlagueTypeController crea una nueva instancia de PlagueTypeController
func NewPlagueTypeController(service *services.PlagueTypeService) *PlagueTypeController {
	return &PlagueTypeController{PlagueTypeService: service}
}

// CrearTipoPlaga maneja la creación de un nuevo tipo de plaga
func (controller *PlagueTypeController) CrearTipoPlaga(c *gin.Context) {
	var tipoPlaga models.PlagueType
	if err := c.ShouldBindJSON(&tipoPlaga); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := controller.PlagueTypeService.CrearTipoPlaga(&tipoPlaga); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tipoPlaga)
}

// ObtenerTodosTiposPlaga maneja la obtención de todos los tipos de plaga
func (controller *PlagueTypeController) ObtenerTodosTiposPlaga(c *gin.Context) {
	tiposPlaga, err := controller.PlagueTypeService.ObtenerTodosTiposPlaga()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los tipos de plaga"})
		return
	}

	c.JSON(http.StatusOK, tiposPlaga)
}

// ObtenerTipoPlagaPorID maneja la obtención de un tipo de plaga por su ID
func (controller *PlagueTypeController) ObtenerTipoPlagaPorID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	tipoPlaga, err := controller.PlagueTypeService.ObtenerTipoPlagaPorID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tipoPlaga)
}

// ActualizarTipoPlaga maneja la actualización de un tipo de plaga
func (controller *PlagueTypeController) ActualizarTipoPlaga(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var updatedTipoPlaga models.PlagueType
	if err := c.ShouldBindJSON(&updatedTipoPlaga); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tipoPlaga, err := controller.PlagueTypeService.ActualizarTipoPlaga(uint(id), &updatedTipoPlaga)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tipoPlaga)
}

// EliminarTipoPlaga maneja la eliminación de un tipo de plaga
func (controller *PlagueTypeController) EliminarTipoPlaga(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := controller.PlagueTypeService.EliminarTipoPlaga(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tipo de plaga eliminado correctamente"})
}
