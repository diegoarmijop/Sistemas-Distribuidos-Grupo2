package services

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sensor-dron-nodo1/models"
	"strconv"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type SensorService struct {
	DB       *gorm.DB
	RabbitMQ *amqp.Channel
}

func NewSensorService(db *gorm.DB, rabbitMQ *amqp.Channel) *SensorService {
	return &SensorService{DB: db, RabbitMQ: rabbitMQ}
}

// CrearSensor crea un nuevo sensor
func (service *SensorService) CrearSensor(sensor *models.Sensor) error {
	return service.DB.Create(sensor).Error
}

// ObtenerTodosSensores obtiene todos los sensores
func (service *SensorService) ObtenerTodosSensores() ([]models.Sensor, error) {
	var sensores []models.Sensor
	err := service.DB.Find(&sensores).Error
	return sensores, err
}

// ObtenerSensorPorID obtiene un sensor por su ID
func (service *SensorService) ObtenerSensorPorID(id string, sensor *models.Sensor) error {
	return service.DB.First(sensor, "id = ?", id).Error
}

// ActualizarSensor actualiza un sensor existente
func (service *SensorService) ActualizarSensor(sensor *models.Sensor) error {
	return service.DB.Save(sensor).Error
}

// EliminarSensor elimina un sensor por su ID
func (service *SensorService) EliminarSensor(id string) error {
	return service.DB.Delete(&models.Sensor{}, "id = ?", id).Error
}

// PublicarDatos publica datos simulados de un sensor en RabbitMQ
func (service *SensorService) PublicarDatos(sensor models.Sensor, sensorID string) error {
	queueName := "sensor." + sensorID
	_, err := service.RabbitMQ.QueueDeclare(
		queueName, true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(sensor)
	if err != nil {
		return err
	}

	return service.RabbitMQ.Publish(
		"", queueName, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// GenerarDatos genera datos simulados para un sensor
func (service *SensorService) GenerarDatos() models.Sensor {
	return models.Sensor{
		Temperatura: service.randomValue(15, 45) + "°C",
		Humedad:     service.randomValue(10, 100) + "%",
		Insectos:    service.randomValue(0, 10),
		Luz:         service.randomValue(100, 1000) + " lux",
	}
}

func (service *SensorService) randomValue(min, max int) string {
	return strconv.Itoa(rand.Intn(max-min) + min)
}

func (service *SensorService) PublicarDatosSensor(sensor models.Sensor, sensorID string) error {
	queueName := "sensor." + sensorID

	// Declarar la cola si no existe
	_, err := service.RabbitMQ.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("error declarando cola del sensor: %v", err)
	}

	// Serializar los datos del sensor
	body, err := json.Marshal(sensor)
	if err != nil {
		return fmt.Errorf("error serializando datos del sensor: %v", err)
	}

	// Publicar los datos en la cola del sensor
	err = service.RabbitMQ.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("error publicando mensaje: %v", err)
	}

	log.Printf("Datos del sensor %s publicados en la cola %s: %+v", sensorID, queueName, sensor)
	return nil
}
