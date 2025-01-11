package controllers

import (
	"go-backend/models"
	"go-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PestTypeController struct {
	service *services.PestTypeService
}

func NewPestTypeController(service *services.PestTypeService) *PestTypeController {
	return &PestTypeController{service: service}
}

func (c *PestTypeController) Create(ctx *gin.Context) {
	var pestType models.PestType
	if err := ctx.ShouldBindJSON(&pestType); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Create(&pestType); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, pestType)
}

func (c *PestTypeController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	pestType, err := c.service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tipo de plaga no encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, pestType)
}

func (c *PestTypeController) GetAll(ctx *gin.Context) {
	pestTypes, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pestTypes)
}

func (c *PestTypeController) Update(ctx *gin.Context) {
	var pestType models.PestType
	if err := ctx.ShouldBindJSON(&pestType); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Update(&pestType); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pestType)
}

func (c *PestTypeController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tipo de plaga eliminado exitosamente"})
}
