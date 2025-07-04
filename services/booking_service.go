package services

import (
	"encoding/json"
	"flyhorizons-bookingservice/config"
	"flyhorizons-bookingservice/models"
	"flyhorizons-bookingservice/models/enums"
	"flyhorizons-bookingservice/services/converter"
	"flyhorizons-bookingservice/services/errors"
	"flyhorizons-bookingservice/services/interfaces"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type BookingService struct {
	bookingRepo        interfaces.BookingRepository
	bookingConverter   converter.BookingConverter
	passengerConverter converter.PassengerConverter
	seatConverter      converter.SeatConverter
}

func NewBookingService(repo interfaces.BookingRepository, bookingConverter converter.BookingConverter, passengerConverter converter.PassengerConverter, seatConverter converter.SeatConverter) *BookingService {
	return &BookingService{
		bookingRepo:        repo,
		bookingConverter:   bookingConverter,
		passengerConverter: passengerConverter,
		seatConverter:      seatConverter,
	}
}

func (s *BookingService) BookingExists(bookingID int) bool {
	for _, booking := range s.bookingRepo.GetAll() {
		if booking.ID == bookingID {
			return true
		}
	}
	return false
}

func (s *BookingService) GetByID(id int) models.Booking {
	bookingEntity := s.bookingRepo.GetByID(id)
	booking := s.bookingConverter.ConvertBookingEntityToBooking(bookingEntity)
	return booking
}

func (s *BookingService) GetByUserID(userID int) []models.Booking {
	bookingEntities := s.bookingRepo.GetByUserID(userID)
	var bookings []models.Booking
	for _, entity := range bookingEntities {
		bookings = append(bookings, s.bookingConverter.ConvertBookingEntityToBooking(entity))
	}
	return bookings
}

func (s *BookingService) Create(booking models.Booking) (*models.Booking, error) {
	if s.BookingExists(booking.ID) {
		return nil, errors.NewBookingExistsError(booking.ID, 409)
	}

	// Set the initial booking status to "Pending"
	// This is when the booking payment has not been (successfully) processed yet
	booking.Status = enums.Pending
	bookingEntity := s.bookingConverter.ConvertBookingToBookingEntity(booking)

	createdEntityPtr := s.bookingRepo.Create(bookingEntity)
	if createdEntityPtr == nil {
		return nil, errors.NewBookingCreateError(booking.ID, 500)
	}

	createdEntity := *createdEntityPtr
	createdBooking := s.bookingConverter.ConvertBookingEntityToBooking(createdEntity)

	// Extract the payment information from the original request
	paymentRequest := models.PaymentRequest{
		BookingID: createdBooking.ID,
		Payment:   booking.Payment,
	}

	// Publish to RabbitMQ
	channel := config.RabbitMQClient.Channel
	body, err := json.Marshal(paymentRequest)
	if err != nil {
		log.Printf("Error marshaling payment details to JSON: %v\n", err)
	} else {
		err = channel.Publish(
			"",
			"booking.created",
			false,
			false,
			amqp091.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			log.Printf("Error publishing booking created event to RabbitMQ: %v\n", err)
		}
	}

	return &createdBooking, nil
}

func (s *BookingService) DeleteByBookingID(id int) (bool, error) {
	if !s.BookingExists(id) {
		return false, errors.NewBookingNotFoundError(id, 404)
	}
	return s.bookingRepo.DeleteByBookingID(id), nil
}

func (s *BookingService) UpdateStatus(bookingID int, status enums.Status) {
	s.bookingRepo.UpdateStatus(bookingID, status)
}

func (s *BookingService) Update(booking models.Booking) (*models.Booking, error) {
	if !s.BookingExists(booking.ID) {
		return nil, errors.NewBookingNotFoundError(booking.ID, 404)
	}

	entity := s.bookingConverter.ConvertBookingToBookingEntity(booking)
	updatedEntity := s.bookingRepo.Update(entity)
	updatedBooking := s.bookingConverter.ConvertBookingEntityToBooking(updatedEntity)

	return &updatedBooking, nil
}
