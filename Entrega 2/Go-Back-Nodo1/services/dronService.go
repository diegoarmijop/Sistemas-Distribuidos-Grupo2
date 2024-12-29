package services

import (
	"encoding/json"
	"log"
	"sensor-dron-nodo1/models"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type DronService struct {
	DB       *gorm.DB
	RabbitMQ *amqp.Channel
}

func NewDronService(db *gorm.DB, rabbitMQ *amqp.Channel) *DronService {
	return &DronService{DB: db, RabbitMQ: rabbitMQ}
}

// CrearDron crea un nuevo dron
func (service *DronService) CrearDron(dron *models.Dron) error {
	return service.DB.Create(dron).Error
}

// ObtenerTodosDrones obtiene todos los drones
func (service *DronService) ObtenerTodosDrones() ([]models.Dron, error) {
	var drones []models.Dron
	err := service.DB.Find(&drones).Error
	return drones, err
}

// ObtenerDronPorID obtiene un dron por su ID
func (service *DronService) ObtenerDronPorID(id string, dron *models.Dron) error {
	return service.DB.First(dron, "id = ?", id).Error
}

// ActualizarDron actualiza un dron existente
func (service *DronService) ActualizarDron(dron *models.Dron) error {
	return service.DB.Save(dron).Error
}

// EliminarDron elimina un dron por su ID
func (service *DronService) EliminarDron(id string) error {
	return service.DB.Delete(&models.Dron{}, "id = ?", id).Error
}

// ProcesarDatosSensor consume datos de sensores y los envía al nodo local
func (service *DronService) ProcesarDatosSensor(sensorID string, dronID string) {
	queueName := "sensor." + sensorID

	// Declarar la cola si no existe
	_, err := service.RabbitMQ.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Printf("Error declarando cola: %v", err)
		return
	}

	msgs, err := service.RabbitMQ.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Printf("Error consumiendo mensajes: %v", err)
		return
	}

	log.Printf("Dron %s esperando mensajes del sensor %s...", dronID, sensorID)

	go func() {
		for msg := range msgs {
			var sensor models.Sensor
			if err := json.Unmarshal(msg.Body, &sensor); err != nil {
				log.Printf("Error procesando mensaje: %v", err)
				continue
			}

			log.Printf("Dron %s recibió datos: %+v", dronID, sensor)

			// Reenviar al nodo local
			service.EnviarANodoLocal(sensor, dronID)
		}
	}()
}

// EnviarANodoLocal publica datos en la cola del nodo local
func (service *DronService) EnviarANodoLocal(sensor models.Sensor, dronID string) {
	queueName := "nodo.local"
	body, err := json.Marshal(map[string]interface{}{
		"dron_id": dronID,
		"sensor":  sensor,
	})
	if err != nil {
		log.Printf("Error serializando mensaje: %v", err)
		return
	}

	err = service.RabbitMQ.Publish(
		"", queueName, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Error publicando mensaje: %v", err)
	}
}
