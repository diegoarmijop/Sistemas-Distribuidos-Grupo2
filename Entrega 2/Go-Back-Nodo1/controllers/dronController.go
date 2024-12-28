package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sensor-dron-nodo1/models"
	"sensor-dron-nodo1/services"
)

type DronController struct {
	DronService *services.DronService
}

func NewDronController(service *services.DronService) *DronController {
	return &DronController{DronService: service}
}

// CrearDron maneja la creación de un nuevo dron
func (controller *DronController) CrearDron(c *gin.Context) {
	var dron models.Dron
	if err := c.ShouldBindJSON(&dron); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Llama al servicio para crear el dron
	if err := controller.DronService.CrearDron(&dron); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el dron: " + err.Error()})
		return
	}

	// Si RutaID no es nulo, cargar RutaActual
	//if dron.RutaID != nil {
	//	if err := controller.DronService.CargarRutaActual(&dron); err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al cargar la ruta asociada: " + err.Error()})
	//		return
	//	}
	//}

	c.JSON(http.StatusOK, dron)
}

// ObtenerTodosDrones devuelve todos los drones
func (controller *DronController) ObtenerTodosDrones(c *gin.Context) {
	drones, err := controller.DronService.ObtenerTodosDrones()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener los drones"})
		return
	}
	c.JSON(http.StatusOK, drones)
}

// ActualizarDron actualiza los datos de un dron existente
func (controller *DronController) ActualizarDron(c *gin.Context) {
	var dron models.Dron
	id := c.Param("id")

	// Verificar que el dron exista
	if err := controller.DronService.ObtenerDronPorID(id, &dron); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dron no encontrado"})
		return
	}

	// Actualizar datos del dron
	if err := c.ShouldBindJSON(&dron); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	if err := controller.DronService.ActualizarDron(&dron); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el dron: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dron)
}

// EliminarDron elimina un dron por su ID
func (controller *DronController) EliminarDron(c *gin.Context) {
	id := c.Param("id")

	if err := controller.DronService.EliminarDron(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el dron: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dron eliminado exitosamente"})
}
