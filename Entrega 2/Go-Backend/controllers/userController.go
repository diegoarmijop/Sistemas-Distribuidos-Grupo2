package controller

import (
	"encoding/json"
	"go-backend/models"
	"go-backend/services"
	"net/http"
	"strconv"
	"strings"
)

// UsuarioController maneja las solicitudes HTTP relacionadas con los usuarios
type UsuarioController struct {
	Service *services.UserServices // Servicio que maneja la lógica de negocio de los usuarios
}

// Nueva instancia del controlador
func NewUsuarioController(service *services.UserServices) *UsuarioController {
	return &UsuarioController{Service: service}
}

// Obtener un usuario por ID
func (uc *UsuarioController) ObtenerUsuario(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/") // Dividir la URL para obtener el ID
	if len(parts) < 3 {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2]) // Convertir el ID de string a int
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Obtener usuario del servicio
	usuario, err := uc.Service.ObtenerUsuarioPorID(id)
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	// Devolver el usuario en formato JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// Crear un nuevo usuario
func (uc *UsuarioController) CrearUsuario(w http.ResponseWriter, r *http.Request) {
	var usuario models.User
	err := json.NewDecoder(r.Body).Decode(&usuario) // Decodificar el JSON en un objeto Usuario
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Crear usuario en el servicio
	usuarioCreado, err := uc.Service.CrearUsuario(usuario)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Devolver el usuario creado en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usuarioCreado)
}

// Obtener todos los usuarios
func (uc *UsuarioController) ObtenerUsuarios(w http.ResponseWriter, r *http.Request) {
	usuarios, err := uc.Service.ObtenerUsuarios() // Obtener todos los usuarios
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Devolver la lista de usuarios en formato JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarios)
}

// Actualizar un usuario
func (uc *UsuarioController) ActualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2]) // Obtener y convertir el ID
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var usuario models.User
	err = json.NewDecoder(r.Body).Decode(&usuario) // Decodificar el cuerpo en un usuario
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Actualizar el usuario a través del servicio
	usuarioActualizado, err := uc.Service.ActualizarUsuario(id, usuario)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Devolver el usuario actualizado en formato JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarioActualizado)
}

// Eliminar un usuario
func (uc *UsuarioController) EliminarUsuario(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2]) // Obtener el ID desde la URL
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Eliminar el usuario a través del servicio
	err = uc.Service.EliminarUsuario(id)
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	// Confirmar la eliminación con un código 204 (sin contenido)
	w.WriteHeader(http.StatusNoContent)
}
