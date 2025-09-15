package queues

import (
	"encoding/json"
	"fmt"
	"log"
	"search-api/domain/courses"

	"github.com/streadway/amqp"
)

// RabbitConfig define la configuración para conectarse a RabbitMQ
type RabbitConfig struct {
	Host      string
	Port      string
	Username  string
	Password  string
	QueueName string
}

// Rabbit representa una conexión de RabbitMQ
type Rabbit struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

// NewRabbit crea una nueva conexión a RabbitMQ y declara la cola
func NewRabbit(config RabbitConfig) Rabbit {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.Username, config.Password, config.Host, config.Port)
	connection, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("Error al conectar con RabbitMQ: %v", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("Error al crear el canal de RabbitMQ: %v", err)
	}

	queue, err := channel.QueueDeclare(
		config.QueueName,
		false, // No persistente
		false, // No autoeliminada
		false, // No exclusiva
		false, // No espera
		nil,
	)
	if err != nil {
		log.Fatalf("Error al declarar la cola: %v", err)
	}

	return Rabbit{
		connection: connection,
		channel:    channel,
		queue:      queue,
	}
}

// StartConsumer inicia la escucha de mensajes en la cola de RabbitMQ
func (rabbit Rabbit) StartConsumer(handler func(courses.CourseUpdate)) error {
	messages, err := rabbit.channel.Consume(
		rabbit.queue.Name,
		"",
		true,  // Auto-acuse de recibo (auto-acknowledge)
		false, // No exclusivo
		false, // No espera
		false, // No local
		nil,
	)
	if err != nil {
		return fmt.Errorf("error al registrar el consumidor: %v", err)
	}

	go func() {
		for msg := range messages {
			var courseUpdate courses.CourseUpdate
			if err := json.Unmarshal(msg.Body, &courseUpdate); err != nil {
				log.Printf("Error al deserializar el mensaje: %v", err)
				continue
			}

			log.Printf("Mensaje recibido en el consumidor: %+v", courseUpdate)

			// Pasar el mensaje al manejador (handler)
			handler(courseUpdate)
		}
	}()

	return nil
}

// Close cierra la conexión y el canal de RabbitMQ
func (rabbit Rabbit) Close() {
	if err := rabbit.channel.Close(); err != nil {
		log.Printf("Error al cerrar el canal de RabbitMQ: %v", err)
	}
	if err := rabbit.connection.Close(); err != nil {
		log.Printf("Error al cerrar la conexión de RabbitMQ: %v", err)
	}
}
