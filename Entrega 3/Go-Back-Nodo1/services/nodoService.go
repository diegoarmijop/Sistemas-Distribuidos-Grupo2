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
	// Inicializar el mensaje final
	var mensajes []string

	// Procesar Humedad
	humedad, err := strconv.Atoi(sensor.Humedad[:len(sensor.Humedad)-1]) // Convertir "60%" a 60
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
	temperaturaStr := strings.TrimRight(sensor.Temperatura, "°C") // Remover caracteres "°C"
	temperatura, err := strconv.Atoi(temperaturaStr)              // Convertir a entero
	if err != nil {
		log.Printf("Error al procesar la temperatura: %v", err)
		return
	}
	switch {
	case temperatura > 40:
		mensajes = append(mensajes, "Temperatura extremadamente alta: peligro extremo.")
	case temperatura > 36:
		mensajes = append(mensajes, "Temperatura alta: riesgo crítico de calor.")
	case temperatura < 5:
		mensajes = append(mensajes, "Temperatura baja: riesgo de hipotermia.")
	default:
		mensajes = append(mensajes, "Temperatura adecuada.")
	}

	// Procesar Insectos
	var insectosMensaje string
	switch sensor.Insectos {
	case "bajo":
		insectosMensaje = "Nivel bajo de insectos: sin riesgo."
	case "medio":
		insectosMensaje = "Nivel medio de insectos: requiere monitoreo."
	case "alto":
		insectosMensaje = "Nivel alto de insectos: riesgo crítico significativo."
	case "abundancia peligrosa":
		insectosMensaje = "Nivel extremadamente alto de insectos: peligro extremo."
	default:
		insectosMensaje = "Nivel de insectos desconocido."
	}
	mensajes = append(mensajes, insectosMensaje)

	// Procesar Luz
	luz, err := strconv.Atoi(sensor.Luz[:len(sensor.Luz)-2]) // Convertir "10UV" a 10
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

	// Armar el mensaje final
	mensajeFinal := fmt.Sprintf("Evaluación de sensor:\n- %s",
		strings.Join(mensajes, "\n- "))

	// Log del mensaje final
	log.Println(mensajeFinal)

	// Si hay parámetros críticos, enviar alerta
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

// Función auxiliar para determinar si hay condiciones críticas en los mensajes
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
// EnviarAlerta estructura y envía una alerta a la base central
func (service *NodoService) EnviarAlerta(sensor models.Sensor, usuarioID uint, eventoPlagaID *uint) {
	// Inicializar variables para determinar el estado y el tipo de alerta
	var (
		estado          string
		tipoAlertas     []string
		parametrosAltos int
	)

	// Evaluar cada parámetro para determinar criticidad
	// 1. Evaluar Humedad
	humedad, _ := strconv.Atoi(sensor.Humedad[:len(sensor.Humedad)-1])
	if humedad < 20 || humedad > 70 {
		tipoAlertas = append(tipoAlertas, "Humedad baja-alta")
		parametrosAltos++
	}

	// 2. Evaluar Temperatura
	temperatura, _ := strconv.Atoi(sensor.Temperatura[:len(sensor.Temperatura)-2])
	if temperatura > 36 {
		tipoAlertas = append(tipoAlertas, "Temperatura alta")
		parametrosAltos++
	}
	if temperatura < 5 {
		tipoAlertas = append(tipoAlertas, "Temperatura baja")
		parametrosAltos++
	}

	// 3. Evaluar Insectos
	switch sensor.Insectos {
	case "alto", "abundancia peligrosa":
		tipoAlertas = append(tipoAlertas, "Nivel alto de insectos")
		parametrosAltos++
	}

	// 4. Evaluar Luz
	luz, _ := strconv.Atoi(sensor.Luz[:len(sensor.Luz)-2])
	if luz > 11 {
		tipoAlertas = append(tipoAlertas, "Nivel de luz extremadamente alto")
		parametrosAltos++
	} else if luz > 8 {
		tipoAlertas = append(tipoAlertas, "Nivel de luz alto")
		parametrosAltos++
	}

	// Determinar el estado en función de los parámetros críticos
	switch {
	case parametrosAltos == 1:
		estado = "Seria"
	case parametrosAltos == 2:
		estado = "Crítico"
	case parametrosAltos > 2:
		estado = "Extremo Peligro"
	}

	// Crear la descripción de la alerta
	descripcion := fmt.Sprintf(
		"Temperatura: %s, Humedad: %s, Insectos: %s, Luz: %s",
		sensor.Temperatura, sensor.Humedad, sensor.Insectos, sensor.Luz,
	)

	// Mapear la alerta para enviar al backend
	alert := map[string]interface{}{
		"estado":          estado,
		"descripcion":     descripcion,
		"fecha_hora":      time.Now().Format(time.RFC3339),
		"tipo_alerta":     strings.Join(tipoAlertas, "/"),
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

	// Manejar la respuesta del backend
	if resp.StatusCode == http.StatusOK {
		log.Printf("Alerta enviada con éxito. Código HTTP: %d", resp.StatusCode)
	} else {
		log.Printf("Error al enviar alerta. Código HTTP: %d", resp.StatusCode)
	}
}
