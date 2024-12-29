package controllers

import (
	"net/http"
	"sensor-dron-nodo1/models"
	"sensor-dron-nodo1/services"

	"github.com/gin-gonic/gin"
)

type NodoController struct {
	NodoService *services.NodoService
}

func NewNodoController(service *services.NodoService) *NodoController {
	return &NodoController{NodoService: service}
}

// CrearNodo maneja la creación de un nuevo nodo
func (controller *NodoController) CrearNodo(c *gin.Context) {
	var nodo models.Nodo
	if err := c.ShouldBindJSON(&nodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	if err := controller.NodoService.CrearNodo(&nodo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el nodo: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, nodo)
}

// ObtenerTodosNodos devuelve todos los nodos
func (controller *NodoController) ObtenerTodosNodos(c *gin.Context) {
	nodos, err := controller.NodoService.ObtenerTodosNodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener los nodos"})
		return
	}
	c.JSON(http.StatusOK, nodos)
}

// ActualizarNodo actualiza los datos de un nodo existente
func (controller *NodoController) ActualizarNodo(c *gin.Context) {
	var nodo models.Nodo
	id := c.Param("id")

	if err := controller.NodoService.ObtenerNodoPorID(id, &nodo); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nodo no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&nodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	if err := controller.NodoService.ActualizarNodo(&nodo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el nodo: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, nodo)
}

// EliminarNodo elimina un nodo por su ID
func (controller *NodoController) EliminarNodo(c *gin.Context) {
	id := c.Param("id")

	if err := controller.NodoService.EliminarNodo(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el nodo: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nodo eliminado exitosamente"})
}
