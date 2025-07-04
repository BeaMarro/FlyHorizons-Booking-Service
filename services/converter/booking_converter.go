package converter

import (
	"flyhorizons-bookingservice/models"
	"flyhorizons-bookingservice/models/enums"
	entities "flyhorizons-bookingservice/repositories/entity"
	"time"
)

type BookingConverter struct {
	passengerConverter PassengerConverter
	seatConverter      SeatConverter
}

func (bookingConverter *BookingConverter) ConvertBookingEntityToBooking(entity entities.BookingEntity) models.Booking {
	return models.Booking{
		ID:          entity.ID,
		UserID:      entity.UserID,
		FlightCode:  entity.FlightCode,
		FlightClass: enums.FlightClassFromInt(entity.FlightClass),
		Luggage:     enums.LuggageClassesFromJSONString(entity.Luggage),
		Seats:       bookingConverter.seatConverter.ConvertSeatEntitiesToSeats(entity.Seats),
		Passengers:  bookingConverter.passengerConverter.ConvertPassengerEntitiesToPassengers(entity.Passengers),
		Status:      enums.Status(entity.Status),
	}
}

func (bookingConverter *BookingConverter) ConvertBookingToBookingEntity(booking models.Booking) entities.BookingEntity {
	bookingEntity := entities.BookingEntity{
		ID:          booking.ID,
		UserID:      booking.UserID,
		FlightCode:  booking.FlightCode,
		FlightClass: int(booking.FlightClass),
		CreatedAt:   time.Now(),
		Luggage:     enums.JSONStringToLuggageClasses(booking.Luggage),
		Status:      string(booking.Status),
	}

	bookingEntity.Passengers = bookingConverter.passengerConverter.ConvertPassengersToPassengerEntities(booking.Passengers, bookingEntity.ID)
	bookingEntity.Seats = bookingConverter.seatConverter.ConvertSeatsToSeatEntities(booking.Seats, bookingEntity.ID)

	return bookingEntity
}
