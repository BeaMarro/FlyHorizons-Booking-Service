package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
)

func InitializeRabbitMQ() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	connection, err := amqp091.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("An error occurred while connecting to RabbitMQ: %s", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("An error occurred while opening the RabbitMQ channel: %s", err)
	}

	// Declare the queues
	queues := []string{
		"booking.created",
		"booking.confirmed",
		"user_deleted",
		"payment.success",
		"payment.failed",
	}

	for _, queueName := range queues {
		_, err = channel.QueueDeclare(
			queueName, // Queue Name
			true,      // Durable
			false,     // Auto Delete
			false,     // Exclusive
			false,     // No Wait
			nil,       // Arguments
		)
		if err != nil {
			log.Fatalf("An error occurred while declaring the '%s' queue: %s", queueName, err)
		}
	}

	RabbitMQClient = &RabbitMQ{
		Connection: connection,
		Channel:    channel,
	}

	log.Println("RabbitMQ has been initialized successfully.")
}
