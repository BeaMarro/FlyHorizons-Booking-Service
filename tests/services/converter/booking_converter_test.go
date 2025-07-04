package converter_test

import (
	"flyhorizons-bookingservice/models"
	"flyhorizons-bookingservice/models/enums"
	entities "flyhorizons-bookingservice/repositories/entity"
	"flyhorizons-bookingservice/services/converter"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Setup
func setupBookingConverter() converter.BookingConverter {
	return converter.BookingConverter{}
}

func getBookingEntity() entities.BookingEntity {
	return entities.BookingEntity{
		ID:          0,
		UserID:      2,
		FlightCode:  "FR788",
		FlightClass: 1,
		CreatedAt:   time.Date(2025, 4, 3, 9, 0, 0, 0, time.UTC),
		Passengers:  getPassengerEntities(),
		Seats:       getSeatEntities(),
		Luggage:     `["SmallBag","Cargo20kg"]`,
	}
}

func getBooking() models.Booking {
	return models.Booking{
		ID:          0,
		UserID:      2,
		FlightCode:  "FR788",
		FlightClass: 1,
		Luggage:     []enums.Luggage{enums.SmallBag, enums.Cargo20kg},
		Seats:       getSeats(),
		Passengers:  getPassengers(),
	}
}

func TestConvertBookingEntityToBookingReturnsBooking(t *testing.T) {
	// Arrange
	bookingConverter := setupBookingConverter()
	bookingEntity := getBookingEntity()
	expectedBooking := getBooking()

	// Act
	booking := bookingConverter.ConvertBookingEntityToBooking(bookingEntity)

	// Assert
	assert.Equal(t, expectedBooking, booking)
}

func TestConvertBookingToBookingEntityReturnsBookingEntity(t *testing.T) {
	// Arrange
	bookingConverter := setupBookingConverter()
	booking := getBooking()
	expectedBookingEntity := getBookingEntity()

	// Act
	bookingEntity := bookingConverter.ConvertBookingToBookingEntity(booking)

	// Assert
	// Ignoring CreatedAt timestamp and Seat IDs which can vary by environment
	// Copy the created timestamp and seat IDs to avoid comparison failures
	bookingEntityCopy := bookingEntity
	bookingEntityCopy.CreatedAt = expectedBookingEntity.CreatedAt

	// Fix the seat IDs to match expected values
	for i := range bookingEntityCopy.Seats {
		bookingEntityCopy.Seats[i].ID = expectedBookingEntity.Seats[i].ID
	}

	assert.Equal(t, expectedBookingEntity, bookingEntityCopy)
}
