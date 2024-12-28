package services

import (
	"errors"
	"go-backend/models"

	"gorm.io/gorm"
)

// PlagueTypeService gestiona las operaciones sobre los tipos de plaga
type PlagueTypeService struct {
	DB *gorm.DB
}

// NewPlagueTypeService crea una nueva instancia de PlagueTypeService
func NewPlagueTypeService(db *gorm.DB) *PlagueTypeService {
	return &PlagueTypeService{DB: db}
}

// CrearTipoPlaga crea un nuevo tipo de plaga si no existe un tipo con el mismo nombre científico
func (service *PlagueTypeService) CrearTipoPlaga(tipoPlaga *models.PlagueType) error {
	// Verificar si ya existe un tipo de plaga con el mismo nombre científico
	var existingPlagueType models.PlagueType
	if err := service.DB.Where("nombre_cientifico = ?", tipoPlaga.NombreCientifico).First(&existingPlagueType).Error; err == nil {
		return errors.New("ya existe un tipo de plaga con este nombre científico")
	}

	// Guardar el nuevo tipo de plaga en la base de datos
	if err := service.DB.Create(&tipoPlaga).Error; err != nil {
		return err
	}
	return nil
}

// ObtenerTodosTiposPlaga obtiene todos los tipos de plaga
func (service *PlagueTypeService) ObtenerTodosTiposPlaga() ([]models.PlagueType, error) {
	var tiposPlaga []models.PlagueType
	if err := service.DB.Find(&tiposPlaga).Error; err != nil {
		return nil, err
	}
	return tiposPlaga, nil
}

// ObtenerTipoPlagaPorID obtiene un tipo de plaga por su ID
func (service *PlagueTypeService) ObtenerTipoPlagaPorID(id uint) (*models.PlagueType, error) {
	var tipoPlaga models.PlagueType
	if err := service.DB.First(&tipoPlaga, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tipo de plaga no encontrado")
		}
		return nil, err
	}
	return &tipoPlaga, nil
}

// ActualizarTipoPlaga actualiza los detalles de un tipo de plaga
func (service *PlagueTypeService) ActualizarTipoPlaga(id uint, updatedPlagueType *models.PlagueType) (*models.PlagueType, error) {
	var tipoPlaga models.PlagueType
	if err := service.DB.First(&tipoPlaga, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tipo de plaga no encontrado")
		}
		return nil, err
	}

	// Actualizar los campos
	tipoPlaga.NombreComun = updatedPlagueType.NombreComun
	tipoPlaga.Descripcion = updatedPlagueType.Descripcion
	tipoPlaga.NombreCientifico = updatedPlagueType.NombreCientifico

	if err := service.DB.Save(&tipoPlaga).Error; err != nil {
		return nil, err
	}

	return &tipoPlaga, nil
}

// EliminarTipoPlaga elimina un tipo de plaga por su ID
func (service *PlagueTypeService) EliminarTipoPlaga(id uint) error {
	if err := service.DB.Delete(&models.PlagueType{}, id).Error; err != nil {
		return err
	}
	return nil
}
