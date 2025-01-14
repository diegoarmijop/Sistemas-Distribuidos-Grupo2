package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sensor-dron-nodo1/models"
	"strconv"
	"strings"
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
	dronIDUint, err := strconv.ParseUint(dronID, 10, 32)
	if err != nil {
		return fmt.Errorf("error al convertir dronID a uint: %v", err)
	}

	nuevaRuta := models.Ruta{
		FechaHoraInicio:  time.Now(),
		FechaHoraTermino: time.Now().Add(1 * time.Hour),
		FlagDron:         fmt.Sprintf("%d", dronIDUint),
	}

	if err := service.RutaService.CrearRuta(&nuevaRuta); err != nil {
		return fmt.Errorf("error creando la ruta: %v", err)
	}

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
	var mensajes []string

	// Procesar Humedad
	humedad, err := strconv.Atoi(strings.TrimRight(sensor.Humedad, "%"))
	if err != nil {
		log.Printf("Error al procesar la humedad: %v", err)
		return
	}

	switch {
	case humedad < 20:
		mensajes = append(mensajes, "Humedad extremadamente baja: riesgo crítico.")
		log.Println("Humedad baja detectada, generando una ruta para el dron ...", dronID)
		if err := service.CrearRuta(sensor, dronID); err != nil {
			log.Printf("Error generando la ruta: %v", err)
		}
	case humedad >= 20 && humedad <= 40:
		mensajes = append(mensajes, "Humedad baja: requiere monitoreo.")
	case humedad > 40 && humedad <= 70:
		mensajes = append(mensajes, "Humedad adecuada.")
	default:
		mensajes = append(mensajes, "Humedad alta: riesgo crítico de saturación.")
	}

	// Procesar Temperatura
	temperaturaStr := strings.TrimRight(sensor.Temperatura, "°C")
	temperatura, err := strconv.Atoi(temperaturaStr)
	if err != nil {
		log.Printf("Error al procesar la temperatura: %v", err)
		return
	}

	switch {
	case temperatura > 30:
		mensajes = append(mensajes, "Temperatura extremadamente alta: peligro extremo.")
	case temperatura > 25:
		mensajes = append(mensajes, "Temperatura alta: riesgo crítico de calor.")
	case temperatura < 15:
		mensajes = append(mensajes, "Temperatura baja: riesgo de hipotermia.")
	default:
		mensajes = append(mensajes, "Temperatura adecuada.")
	}

	// Procesar Insectos
	switch strings.ToLower(sensor.Insectos) {
	case "bajo":
		mensajes = append(mensajes, "Nivel bajo de insectos: sin riesgo.")
	case "medio":
		mensajes = append(mensajes, "Nivel medio de insectos: requiere monitoreo.")
	case "alto":
		mensajes = append(mensajes, "Nivel alto de insectos: riesgo crítico significativo.")
	case "abundancia peligrosa":
		mensajes = append(mensajes, "Nivel extremadamente alto de insectos: peligro extremo.")
	default:
		mensajes = append(mensajes, "Nivel de insectos desconocido.")
	}

	// Procesar Luz
	luzStr := strings.TrimRight(sensor.Luz, "UV")
	luz, err := strconv.Atoi(luzStr)
	if err != nil {
		log.Printf("Error al procesar la luz: %v", err)
		return
	}

	switch {
	case luz > 11:
		mensajes = append(mensajes, "Nivel de luz extremadamente alto: peligro extremo.")
	case luz > 8:
		mensajes = append(mensajes, "Nivel de luz alto: riesgo crítico de exposición.")
	case luz < 3:
		mensajes = append(mensajes, "Nivel de luz bajo: posible riesgo de fotosíntesis insuficiente.")
	default:
		mensajes = append(mensajes, "Nivel de luz adecuado.")
	}

	mensajeFinal := fmt.Sprintf("Evaluación de sensor:\n- %s",
		strings.Join(mensajes, "\n- "))

	log.Println(mensajeFinal)

	if containsCriticalConditions(mensajes) {
		log.Println("Se detectaron condiciones críticas, enviando alerta...")
		service.EnviarAlerta(sensor, 1, nil)
	} else {
		log.Println("Se detectaron parámetros medianamente alterados, generando una ruta de solución para el dron ...", dronID)
		if err := service.CrearRuta(sensor, dronID); err != nil {
			log.Printf("Error generando la ruta: %v", err)
		}
	}
}

func containsCriticalConditions(mensajes []string) bool {
	for _, mensaje := range mensajes {
		if strings.Contains(mensaje, "riesgo crítico") || strings.Contains(mensaje, "peligro extremo") {
			return true
		}
	}
	return false
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

		dronID, ok := payload["dron_id"].(string)
		if !ok || dronID == "" {
			log.Printf("Error: dron_id no está presente o no es válido")
			continue
		}

		sensorData, ok := payload["sensor"].(map[string]interface{})
		if !ok || sensorData == nil {
			log.Printf("Error: sensor no está presente o no es válido")
			continue
		}

		sensor := models.Sensor{
			Temperatura: sensorData["temperatura"].(string),
			Humedad:     sensorData["humedad"].(string),
			Insectos:    sensorData["insectos"].(string),
			Luz:         sensorData["luz"].(string),
		}

		service.ProcesarSensor(sensor, dronID)
	}
}

func (service *NodoService) EnviarAlerta(sensor models.Sensor, usuarioID uint, eventoPlagaID *uint) {
	var (
		estado          string
		tipoAlertas     []string
		parametrosAltos int
	)

	// Evaluar Humedad
	humedad, _ := strconv.Atoi(strings.TrimRight(sensor.Humedad, "%"))
	if humedad < 20 || humedad > 70 {
		tipoAlertas = append(tipoAlertas, "Humedad baja-alta")
		parametrosAltos++
	}

	// Evaluar Temperatura
	temperaturaStr := strings.TrimRight(sensor.Temperatura, "°C")
	temperatura, err := strconv.Atoi(temperaturaStr)
	if err != nil {
		log.Printf("Error al procesar temperatura: %v", err)
		return
	}

	if temperatura > 30 {
		tipoAlertas = append(tipoAlertas, "Temperatura alta")
		parametrosAltos++
	} else if temperatura < 15 {
		tipoAlertas = append(tipoAlertas, "Temperatura baja")
		parametrosAltos++
	}

	// Evaluar Insectos
	switch strings.ToLower(sensor.Insectos) {
	case "alto", "abundancia peligrosa":
		tipoAlertas = append(tipoAlertas, "Nivel alto de insectos")
		parametrosAltos++
	}

	// Evaluar Luz
	luzStr := strings.TrimRight(sensor.Luz, "UV")
	luz, err := strconv.Atoi(luzStr)
	if err != nil {
		log.Printf("Error al procesar luz: %v", err)
		return
	}

	if luz > 11 {
		tipoAlertas = append(tipoAlertas, "Nivel de luz extremadamente alto")
		parametrosAltos++
	} else if luz > 8 {
		tipoAlertas = append(tipoAlertas, "Nivel de luz alto")
		parametrosAltos++
	}

	// Determinar estado
	switch {
	case parametrosAltos == 0:
		estado = "Normal"
	case parametrosAltos == 1:
		estado = "Seria"
	case parametrosAltos == 2:
		estado = "Crítico"
	case parametrosAltos > 2:
		estado = "Extremo Peligro"
	}

	// Solo enviar alerta si hay parámetros altos
	if parametrosAltos > 0 {
		descripcion := fmt.Sprintf(
			"Temperatura: %s, Humedad: %s, Insectos: %s, Luz: %s",
			sensor.Temperatura, sensor.Humedad, sensor.Insectos, sensor.Luz,
		)

		alert := map[string]interface{}{
			"estado":          estado,
			"descripcion":     descripcion,
			"fecha_hora":      time.Now().Format(time.RFC3339),
			"tipo_alerta":     strings.Join(tipoAlertas, "/"),
			"usuario_id":      usuarioID,
			"evento_plaga_id": eventoPlagaID,
		}

		body, err := json.Marshal(alert)
		if err != nil {
			log.Printf("Error serializando alerta: %v", err)
			return
		}

		resp, err := http.Post(service.BaseCentral+"/api/alertas/", "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Printf("Error enviando alerta: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			log.Printf("Alerta enviada con éxito. Estado: %s, Tipos: %v", estado, tipoAlertas)
		} else {
			log.Printf("Error al enviar alerta. Código HTTP: %d", resp.StatusCode)
		}
	}
}
