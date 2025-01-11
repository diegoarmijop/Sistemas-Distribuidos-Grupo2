package config

import (
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

type RabbitMQConfig struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func InitRabbitMQ() *RabbitMQConfig {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		log.Fatal("La URL de RabbitMQ (RABBITMQ_URL) no está definida en las variables de entorno")
	}

	log.Printf("Intentando conectar a RabbitMQ en: %s", rabbitMQURL)

	// Intentar conexión con reintentos
	var conn *amqp.Connection
	var err error
	for i := 0; i < 3; i++ {
		conn, err = amqp.Dial(rabbitMQURL)
		if err == nil {
			break
		}
		log.Printf("Intento %d: Error conectando a RabbitMQ: %v", i+1, err)
		time.Sleep(time.Second * 2)
	}

	if err != nil {
		log.Printf("Error final conectando a RabbitMQ: %v", err)
		return nil
	}

	log.Println("Conexión a RabbitMQ establecida exitosamente")

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Error abriendo canal en RabbitMQ: %v", err)
		conn.Close()
		return nil
	}

	// Declarar las colas necesarias
	_, err = ch.QueueDeclare(
		"nodo.local", // nombre
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Printf("Error declarando cola del nodo local: %v", err)
		return nil
	}

	log.Println("Canal de RabbitMQ creado y colas declaradas exitosamente")

	return &RabbitMQConfig{
		Connection: conn,
		Channel:    ch,
	}
}
