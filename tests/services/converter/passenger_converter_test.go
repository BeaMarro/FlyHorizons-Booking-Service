package converter_test

import (
	"flyhorizons-bookingservice/models"
	entities "flyhorizons-bookingservice/repositories/entity"
	"flyhorizons-bookingservice/services/converter"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Setup
func setupPassengerConverter() converter.PassengerConverter {
	return converter.PassengerConverter{}
}

func getPassengers() []models.Passenger {
	return []models.Passenger{
		{
			ID:             1,
			FullName:       "John Doe",
			DateOfBirth:    time.Date(1985, 7, 9, 1, 0, 0, 0, time.UTC),
			PassportNumber: "1234",
		},
		{
			ID:             2,
			FullName:       "Jane Doe",
			DateOfBirth:    time.Date(1986, 8, 8, 2, 30, 0, 0, time.UTC),
			PassportNumber: "4321",
		},
	}
}

func getBookingID() int {
	return 123 // Example booking ID
}

func getExpectedPassengerEntities() []entities.PassengerEntity {
	bookingID := getBookingID()
	return []entities.PassengerEntity{
		{
			ID:             1,
			BookingID:      bookingID,
			FullName:       "John Doe",
			DateOfBirth:    time.Date(1985, 7, 9, 1, 0, 0, 0, time.UTC),
			PassportNumber: "1234",
		},
		{
			ID:             2,
			BookingID:      bookingID,
			FullName:       "Jane Doe",
			DateOfBirth:    time.Date(1986, 8, 8, 2, 30, 0, 0, time.UTC),
			PassportNumber: "4321",
		},
	}
}

func getPassengerEntities() []entities.PassengerEntity {
	return []entities.PassengerEntity{
		{
			ID:             1,
			BookingID:      0,
			FullName:       "John Doe",
			DateOfBirth:    time.Date(1985, 7, 9, 1, 0, 0, 0, time.UTC),
			PassportNumber: "1234",
		},
		{
			ID:             2,
			BookingID:      0,
			FullName:       "Jane Doe",
			DateOfBirth:    time.Date(1986, 8, 8, 2, 30, 0, 0, time.UTC),
			PassportNumber: "4321",
		},
	}
}

func TestConvertPassengersToPassengerEntities(t *testing.T) {
	// Arrange
	passengerConverter := setupPassengerConverter()
	passengers := getPassengers()
	bookingID := getBookingID()
	expectedEntities := getExpectedPassengerEntities()

	// Act
	passengerEntities := passengerConverter.ConvertPassengersToPassengerEntities(passengers, bookingID)

	// Assert
	assert.Equal(t, expectedEntities, passengerEntities)
}

func TestConvertPassengerEntitiesToPassengers(t *testing.T) {
	// Arrange
	passengerConverter := setupPassengerConverter()
	passengerEntities := getPassengerEntities()
	expectedPassengers := getPassengers()

	// Act
	passengers := passengerConverter.ConvertPassengerEntitiesToPassengers(passengerEntities)

	// Assert
	assert.Equal(t, len(expectedPassengers), len(passengers))

	assert.Equal(t, expectedPassengers, passengers)
}
