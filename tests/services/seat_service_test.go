package services_test

import (
	entities "flyhorizons-bookingservice/repositories/entity"
	"flyhorizons-bookingservice/services"
	"flyhorizons-bookingservice/services/converter"
	mock_repositories "flyhorizons-bookingservice/tests/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestSeatService struct {
}

// Setup
func setupSeatService() (*mock_repositories.MockSeatRepository, *services.SeatService) {
	mockRepo := new(mock_repositories.MockSeatRepository)
	seatConverter := converter.SeatConverter{}
	seatService := services.NewSeatService(mockRepo, seatConverter)
	return mockRepo, seatService
}

func getSeatOptionEntites() []entities.SeatOptionEntity {
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

// Service Unit Tests
func TestGetSeatsByValidFlightIDReturnsSeats(t *testing.T) {
	// Arrange
	mockRepo, seatService := setupSeatService()
	flightCode := "FR788"
	seatOptionEntities := getSeatOptionEntites()
	expectedSeats := getSeats()
	mockRepo.On("GetByFlightCode", flightCode).Return(seatOptionEntities)

	// Act
	seats, err := seatService.GetByFlightCode(flightCode)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedSeats, seats)
}
