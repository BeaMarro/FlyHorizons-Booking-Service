package main

import (
	"flyhorizons-bookingservice/config"
	"flyhorizons-bookingservice/internal/health"
	"flyhorizons-bookingservice/internal/metrics"
	"flyhorizons-bookingservice/repositories"
	"flyhorizons-bookingservice/routes"
	"flyhorizons-bookingservice/services"
	"flyhorizons-bookingservice/services/authentication"
	"flyhorizons-bookingservice/services/converter"
	"log"

	"github.com/gin-gonic/gin"

	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	healthcfg "github.com/tavsec/gin-healthcheck/config"

	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	// Initialize RabbitMQ for messaging
	config.InitializeRabbitMQ()

	// Setup close for RabbitMQ connections
	if config.RabbitMQClient != nil {
		defer func() {
			if config.RabbitMQClient.Channel != nil {
				err := config.RabbitMQClient.Channel.Close()
				if err != nil {
					log.Printf("Error closing RabbitMQ channel: %v", err)
				}
			}
			if config.RabbitMQClient.Connection != nil {
				err := config.RabbitMQClient.Connection.Close()
				if err != nil {
					log.Printf("Error closing RabbitMQ connection: %v", err)
				}
			}
		}()
	}

	router := gin.Default()

	// Repositories
	baseRepo := repositories.BaseRepository{}
	dbCheck := health.DatabaseCheck{Repository: &baseRepo}

	// -- Health checks setup ---
	conf := healthcfg.DefaultConfig()
	rabbitMQCheck := health.RabbitMQCheck{}
	healthcheck.New(router, conf, []checks.Check{dbCheck, rabbitMQCheck})

	// --- Metrics setup ---
	metrics.RegisterMetricsRoutes(router, dbCheck, rabbitMQCheck)

	// --- Microservices setup ---
	bookingRepo := repositories.NewBookingRepository(&baseRepo)
	seatRepo := repositories.NewSeatRepository(&baseRepo)

	// Converters
	bookingConverter := converter.BookingConverter{}
	passengerConverter := converter.PassengerConverter{}
	seatConverter := converter.SeatConverter{}

	// Authentication
	gatewayAuthMiddleware := authentication.NewGatewayAuthMiddleware()

	// Services
	bookingService := services.NewBookingService(bookingRepo, bookingConverter, passengerConverter, seatConverter)
	seatService := services.NewSeatService(seatRepo, seatConverter)

	// Start the UserEventListener in a goroutine to not block the main thread
	userDeletedListener := services.NewUserEventListener(config.RabbitMQClient, *bookingService)
	go userDeletedListener.StartUserDeletedConsumer()
	log.Println("User deleted consumer started in background")

	// Start the PaymentProcessedListener
	paymentProcessedListener := services.NewPaymentEventListener(config.RabbitMQClient, *bookingService)
	go paymentProcessedListener.StartPaymentProcessedConsumers()
	log.Println("Payment processed consumer started in background")

	// Routes
	routes.RegisterBookingRoutes(router, bookingService, gatewayAuthMiddleware)
	routes.RegisterSeatRoutes(router, seatService)

	// Run the microservice
	log.Println("Starting booking service on port 8083")
	router.Run(":8083")
}
