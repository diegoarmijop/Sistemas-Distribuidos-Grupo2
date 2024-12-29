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
	dronIDUint, err := strconv.ParseUint("8", 10, 32)
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
		Where("id = ?", uint(dronIDUint)). // Usar uint aquí
		Update("ruta_id", nuevaRuta.ID).Error; err != nil {
		return fmt.Errorf("error asignando ruta al dron: %v", err)
	}

	log.Printf("Ruta creada y asignada exitosamente: %+v", nuevaRuta)
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
		log.Println("Humedad baja detectada, generando una ruta...")
		if err := service.CrearRuta(sensor, dronID); err != nil {
			log.Printf("Error generando la ruta: %v", err)
		}
	} else {
		log.Println("Parámetros normales, enviando alerta crítica...")
		service.EnviarAlerta(map[string]interface{}{
			"sensor": sensor,
		})
	}
}

// ProcesarDron consume mensajes de la cola del nodo local
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

		log.Printf("Nodo local recibió: %+v", payload)

		// Procesar datos del sensor
		sensorData := payload["sensor"].(map[string]interface{})
		dronID := payload["dron_id"].(string)

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
func (service *NodoService) EnviarAlerta(data map[string]interface{}) {
	alert := map[string]interface{}{
		"estado":      "Crítico",
		"descripcion": "Parámetros alterados detectados por el nodo.",
		"fecha_hora":  time.Now().Format(time.RFC3339),
		"tipo_alerta": "Sensor Crítico",
		"nodo_id":     "Nodo-1",
		"sensor_data": data,
	}

	body, err := json.Marshal(alert)
	if err != nil {
		log.Printf("Error serializando alerta: %v", err)
		return
	}

	resp, err := http.Post(service.BaseCentral+"/api/alerta", "application/json", bytes.NewBuffer(body))
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
