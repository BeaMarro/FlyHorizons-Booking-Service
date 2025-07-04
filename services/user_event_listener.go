package services

import (
	"encoding/json"
	"flyhorizons-bookingservice/config"
	"flyhorizons-bookingservice/models"
	"log"
)

type UserEventListener struct {
	rabbitMQClient *config.RabbitMQ
	bookingService BookingService
}

func NewUserEventListener(client *config.RabbitMQ, service BookingService) *UserEventListener {
	return &UserEventListener{
		rabbitMQClient: client,
		bookingService: service,
	}
}

func (userEventListener *UserEventListener) StartUserDeletedConsumer() {
	channel := userEventListener.rabbitMQClient.Channel

	messages, err := channel.Consume(
		"user_deleted", // Queue Name
		"",             // Consumer tag
		true,           // Auto acknowledge
		false,          // Exclusive
		false,          // No-local
		false,          // No-wait
		nil,            // Arguments
	)
	if err != nil {
		log.Fatalf("An error occurred while registering the consumer: %v", err)
	}

	// Log consumer setup
	log.Printf("Started consumer for RabbitMQ queue: %s", "user_deleted")
	log.Printf("Channel: RabbitMQ channel is active")

	// Start goroutine to process messages
	go func() {
		for message := range messages {
			// Log the message body
			log.Printf("Received user deletion message: %s", string(message.Body))

			var event models.UserDeletedEvent
			if err := json.Unmarshal(message.Body, &event); err != nil {
				log.Printf("An error occurred while converting the JSON to UserDeletedEvent: %v", err)
				continue
			}

			userID := event.UserID
			// Delete user data
			bookings := userEventListener.bookingService.GetByUserID(userID)
			for _, booking := range bookings {
				userEventListener.bookingService.DeleteByBookingID(booking.ID)
			}
			log.Printf("Successfully deleted the user data for UserID: %d from the booking database", userID)
		}
		// Log when the consumer stops (e.g., if the channel closes)
		log.Printf("Consumer for queue %s stopped", "user_deleted")
	}()

	log.Println("User deletion listener started successfully")
}
