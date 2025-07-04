package converter_test

import (
	"flyhorizons-bookingservice/models"
	entities "flyhorizons-bookingservice/repositories/entity"
	"flyhorizons-bookingservice/services/converter"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Setup
func setupSeatConverter() converter.SeatConverter {
	return converter.SeatConverter{}
}

func getSeats() []models.Seat {
	return []models.Seat{
		{
			Row:       1,
			Column:    "A",
			Available: true,
		},
		{
			Row:       1,
			Column:    "B",
			Available: true,
		},
	}
}

func getExpectedSeatEntities() []entities.SeatEntity {
	bookingID := getBookingID()
	return []entities.SeatEntity{
		{
			BookingID: bookingID,
			Row:       1,
			Column:    "A",
		},
		{
			BookingID: bookingID,
			Row:       1,
			Column:    "B",
		},
	}
}

func getSeatEntities() []entities.SeatEntity {
	return []entities.SeatEntity{
		{
			ID:        1,
			BookingID: 0,
			Row:       1,
			Column:    "A",
		},
		{
			ID:        2,
			BookingID: 0,
			Row:       1,
			Column:    "B",
		},
	}
}

func getSeatOptionEntities() []entities.SeatOptionEntity {
	return []entities.SeatOptionEntity{
		{
			ID:     1,
			Row:    1,
			Column: "A",
			Status: true,
		},
		{
			ID:     2,
			Row:    1,
			Column: "B",
			Status: true,
		},
	}
}

func TestConvertSeatsToSeatEntities(t *testing.T) {
	// Arrange
	seatConverter := setupSeatConverter()
	seats := getSeats()
	bookingID := getBookingID()
	expectedEntities := getExpectedSeatEntities()

	// Act
	seatEntities := seatConverter.ConvertSeatsToSeatEntities(seats, bookingID)

	// Assert
	assert.Equal(t, expectedEntities, seatEntities)
}

func TestConvertSeatEntitiesToSeats(t *testing.T) {
	// Arrange
	seatConverter := setupSeatConverter()
	seatEntities := getSeatEntities()
	expectedSeats := getSeats()

	// Act
	seats := seatConverter.ConvertSeatEntitiesToSeats(seatEntities)

	// Assert
	assert.Equal(t, expectedSeats, seats)
}

func TestConverSeatOptionEntitiesToSeats(t *testing.T) {
	// Arrange
	seatConverter := setupSeatConverter()
	seatOptionEntites := getSeatOptionEntities()
	expectedSeats := getSeats()

	// Act
	seats := seatConverter.ConvertSeatOptionEntitiesToSeats(seatOptionEntites)

	// Assert
	assert.Equal(t, expectedSeats, seats)
}
