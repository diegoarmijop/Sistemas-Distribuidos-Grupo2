package controllers

import (
	"go-backend/models"
	"go-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CampoController maneja las peticiones relacionadas con los campos
type CampoController struct {
	CampoService *services.CampService
}

// NewCampoController crea un nuevo CampoController
func NewCampoController(service *services.CampService) *CampoController {
	return &CampoController{CampoService: service}
}

// CrearCampo maneja la creación de un nuevo campo
func (controller *CampoController) CrearCampo(c *gin.Context) {
	var campo models.Camp
	if err := c.ShouldBindJSON(&campo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Usar el servicio para crear el campo
	if err := controller.CampoService.CrearCampo(&campo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, campo)
}

// ObtenerTodosCampos permite obtener todos los campos
func (controller *CampoController) ObtenerTodosCampos(c *gin.Context) {
	// Usar el servicio para obtener todos los campos
	campos, err := controller.CampoService.ObtenerTodosCampos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener los campos"})
		return
	}

	// Responder con los campos obtenidos
	c.JSON(http.StatusOK, campos)
}

// ObtenerCamposPorUbicacion obtiene los campos por ubicación
func (controller *CampoController) ObtenerCamposPorUbicacion(c *gin.Context) {
	ubicacion := c.Param("ubicacion")

	campos, err := controller.CampoService.ObtenerCamposPorUbicacion(ubicacion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, campos)
}

// ActualizarCultivoCampo permite actualizar el cultivo de un campo
func (controller *CampoController) ActualizarCultivoCampo(c *gin.Context) {
	// Obtener el id de la URL y convertirlo de string a uint
	idStr := c.Param("id")                      // El id se obtiene como un string de la URL
	id, err := strconv.ParseUint(idStr, 10, 32) // Convertir el id a uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Convertir id de uint64 a uint (porque el modelo espera uint)
	campoID := uint(id)

	// Crear un struct para recibir el tipo de cultivo desde el cuerpo de la solicitud
	var input struct {
		TipoCultivo string `json:"tipo_cultivo"`
	}

	// Bind los datos JSON a la estructura
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Llamar al servicio para actualizar el cultivo
	campo, err := controller.CampoService.ActualizarCultivoCampo(campoID, input.TipoCultivo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Responder con el campo actualizado
	c.JSON(http.StatusOK, campo)
}

// controllers/campo_controller.go
func (controller *CampoController) ObtenerResumenCampos(c *gin.Context) {
	resumen, err := controller.CampoService.ObtenerResumenCampos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener el resumen de campos",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resumen)
}
