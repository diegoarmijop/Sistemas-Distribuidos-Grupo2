package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sensor-dron-nodo1/models"
	"strconv"
	"time"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type NodoService struct {
	DB          *gorm.DB
	RabbitMQ    *amqp.Channel
	BaseCentral string
	RutaService *RutaService
}

func NewNodoService(db *gorm.DB, rabbitMQ *amqp.Channel, baseCentral string, rutaService *RutaService) *NodoService {
	return &NodoService{DB: db, RabbitMQ: rabbitMQ, BaseCentral: baseCentral, RutaService: rutaService}
}

// CrearNodo crea un nuevo nodo
func (service *NodoService) CrearNodo(nodo *models.Nodo) error {
	return service.DB.Create(nodo).Error
}

// ObtenerTodosNodos obtiene todos los nodos
func (service *NodoService) ObtenerTodosNodos() ([]models.Nodo, error) {
	var nodos []models.Nodo
	err := service.DB.Find(&nodos).Error
	return nodos, err
}

// ObtenerNodoPorID obtiene un nodo por su ID
func (service *NodoService) ObtenerNodoPorID(id string, nodo *models.Nodo) error {
	return service.DB.First(nodo, "id = ?", id).Error
}

// ActualizarNodo actualiza un nodo existente
func (service *NodoService) ActualizarNodo(nodo *models.Nodo) error {
	return service.DB.Save(nodo).Error
}

// EliminarNodo elimina un nodo por su ID
func (service *NodoService) EliminarNodo(id string) error {
	return service.DB.Delete(&models.Nodo{}, "id = ?", id).Error
}

// CrearRuta genera una nueva ruta y la asigna a un dron
func (service *NodoService) CrearRuta(sensor models.Sensor, dronID string) error {
	// Convertir dronID de string a uint
	dronIDUint, err := strconv.ParseUint(dronID, 10, 32)
	if err != nil {
		return fmt.Errorf("error al convertir dronID a uint: %v", err)
	}

	nuevaRuta := models.Ruta{
		FechaHoraInicio:  time.Now(),
		FechaHoraTermino: time.Now().Add(1 * time.Hour), // Duración estimada de 1 hora
		FlagDron:         fmt.Sprintf("%d", dronIDUint), // Almacenar como string en la ruta
	}

	// Guardar la ruta usando RutaService
	if err := service.RutaService.CrearRuta(&nuevaRuta); err != nil {
		return fmt.Errorf("error creando la ruta: %v", err)
	}

	// Actualizar el dron con la ruta asignada
	if err := service.DB.Model(&models.Dron{}).
		Where("id = ?", uint(dronIDUint)).
		Update("ruta_id", nuevaRuta.ID).Error; err != nil {
		return fmt.Errorf("error asignando ruta al dron: %v", err)
	}

	log.Printf("Ruta creada y asignada exitosamente al dron %d: %+v", dronIDUint, nuevaRuta)
	return nil
}

// ProcesarSensor analiza los datos de un sensor y toma decisiones
func (service *NodoService) ProcesarSensor(sensor models.Sensor, dronID string) {
	humedad, err := strconv.Atoi(sensor.Humedad[:len(sensor.Humedad)-1]) // Convertir "15%" a 15
	if err != nil {
		log.Printf("Error al procesar la humedad: %v", err)
		return
	}

	if humedad < 20 {
		log.Println("Humedad baja detectada, generando una ruta para el dron ...", dronID)
		if err := service.CrearRuta(sensor, dronID); err != nil {
			log.Printf("Error generando la ruta: %v", err)
		}
	} else {
		log.Println("Parámetros críticos, enviando alerta...")
		service.EnviarAlerta(sensor, 1, nil)
	}
}

func (service *NodoService) ProcesarDron() {
	queueName := "nodo.local"
	msgs, err := service.RabbitMQ.Consume(
		queueName, "", true, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("Error consumiendo mensajes: %v", err)
	}

	for msg := range msgs {
		var payload map[string]interface{}
		if err := json.Unmarshal(msg.Body, &payload); err != nil {
			log.Printf("Error procesando mensaje: %v", err)
			continue
		}

		log.Printf("Nodo local recibió datos: %+v", payload)

		// Validar y obtener `dron_id`
		dronID, ok := payload["dron_id"].(string)
		if !ok || dronID == "" {
			log.Printf("Error: dron_id no está presente o no es válido")
			continue
		}

		// Validar y obtener `sensor`
		sensorData, ok := payload["sensor"].(map[string]interface{})
		if !ok || sensorData == nil {
			log.Printf("Error: sensor no está presente o no es válido")
			continue
		}

		// Convertir datos del sensor
		sensor := models.Sensor{
			Temperatura: sensorData["temperatura"].(string),
			Humedad:     sensorData["humedad"].(string),
			Insectos:    sensorData["insectos"].(string),
			Luz:         sensorData["luz"].(string),
		}

		// Llamar a ProcesarSensor
		service.ProcesarSensor(sensor, dronID)
	}
}

// EnviarAlerta estructura y envía una alerta a la base central
func (service *NodoService) EnviarAlerta(sensor models.Sensor, usuarioID uint, eventoPlagaID *uint) {
	// Generar descripción dinámica según los datos del sensor
	descripcion := fmt.Sprintf(
		"Alerta generada por parámetros alterados: Temperatura: %s, Humedad: %s, Insectos: %s, Luz: %s",
		sensor.Temperatura, sensor.Humedad, sensor.Insectos, sensor.Luz,
	)

	alert := map[string]interface{}{
		"estado":          "Crítico",
		"descripcion":     descripcion,
		"fecha_hora":      time.Now().Format(time.RFC3339),
		"tipo_alerta":     "Parámetros Alterados",
		"usuario_id":      usuarioID,
		"evento_plaga_id": eventoPlagaID,
	}

	// Serializar la alerta
	body, err := json.Marshal(alert)
	if err != nil {
		log.Printf("Error serializando alerta: %v", err)
		return
	}

	// Enviar la alerta al backend de la casa central
	resp, err := http.Post(service.BaseCentral+"/api/alertas/", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error enviando alerta: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("Alerta enviada con éxito. Código HTTP: %d", resp.StatusCode)
	} else {
		log.Printf("Error al enviar alerta. Código HTTP: %d", resp.StatusCode)
	}
}
