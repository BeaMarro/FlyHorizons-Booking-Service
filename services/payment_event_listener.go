package services

import (
	"encoding/json"
	"flyhorizons-bookingservice/config"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type PaymentEventListener struct {
	rabbitMQClient *config.RabbitMQ
	bookingService BookingService
}

func NewPaymentEventListener(client *config.RabbitMQ, service BookingService) *PaymentEventListener {
	return &PaymentEventListener{
		rabbitMQClient: client,
		bookingService: service,
	}
}

func (p *PaymentEventListener) StartPaymentProcessedConsumers() {
	channel := p.rabbitMQClient.Channel

	// === Consumer for payment.success ===
	successMessages, err := channel.Consume(
		"payment.success", // Queue name
		"",                // Consumer tag
		true,              // Auto-ack
		false,             // Exclusive
		false,             // No-local
		false,             // No-wait
		nil,               // Args
	)
	if err != nil {
		log.Fatalf("Error consuming payment.success: %v", err)
	}

	go func() {
		for msg := range successMessages {
			log.Printf("[payment.success] Received message: %s", string(msg.Body))

			// Success: Post the booking.confirmed event to the RabbitMQ
			// This is listened by the Email Service, therefore a confirmation email will be sent consecutively

			// --- Publish to RabbitMQ ---
			// Get the bookingID
			var bookingID int
			err := json.Unmarshal(msg.Body, &bookingID)
			if err != nil {
				log.Printf("Error unmarshaling booking ID: %v", err)
				continue // skip this message and continue looping
			}

			log.Printf("BookingID: %v", bookingID)

			// Update the booking status to 'Success'
			p.bookingService.UpdateStatus(bookingID, "Success")

			// Get the booking by bookingID
			booking := p.bookingService.GetByID(bookingID)

			log.Print(booking)

			// Marshal the booking and send it using RabbitMQ
			body, err := json.Marshal(booking)
			if err != nil {
				log.Printf("Error marshaling booking: %v", err)
			} else {
				err = channel.Publish(
					"",
					"booking.confirmed",
					true,
					false,
					amqp091.Publishing{
						ContentType: "application/json",
						Body:        body,
					},
				)
				if err != nil {
					log.Printf("Error publishing booking created event to RabbitMQ: %v\n", err)
				}
				log.Print("Booking information event published to booking.confirmed")
			}
		}
	}()

	// === Consumer for payment.fail ===
	failMessages, err := channel.Consume(
		"payment.failed", // Queue name
		"",               // Consumer tag
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error consuming payment.fail: %v", err)
	}

	go func() {
		for msg := range failMessages {
			log.Printf("[payment.fail] Received message: %s", string(msg.Body))

			// Failed: Deletes the booking from the database
			// The booking is deleted, as the payment was not processed successfully

			// --- Delete Booking by ID ---
			// Get the bookingID
			var bookingID int
			err := json.Unmarshal(msg.Body, &bookingID)
			if err != nil {
				log.Printf("Error unmarshaling booking ID: %v", err)
				continue // skip this message and continue looping
			}

			// Delete booking by bookingID
			p.bookingService.DeleteByBookingID(bookingID)
		}
	}()

	log.Println("Payment event consumers started: payment.success and payment.failed")
}
