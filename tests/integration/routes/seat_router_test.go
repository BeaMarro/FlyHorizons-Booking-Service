package routes_test

import (
	"encoding/json"
	"flyhorizons-bookingservice/models"
	"flyhorizons-bookingservice/routes"
	mock_repositories "flyhorizons-bookingservice/tests/mocks"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestSeatRouter struct {
}

// Setup
func setupSeatRouter(mockService *mock_repositories.MockSeatService) *gin.Engine {
	router := gin.Default()
	routes.RegisterSeatRoutes(router, mockService)
	return router
}

func getFlightSeats() []models.Seat {
	return []models.Seat{
		{
			Row:       1,
			Column:    "A",
			Available: true,
		},
		{
			Row:       1,
			Column:    "B",
			Available: false,
		},
	}
}

func TestGetSeatsByFlightCodeReturnsSeatsJSON(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockSeatService)
	mockSeats := getFlightSeats()
	flightCode := "FR788"
	mockService.On("GetByFlightCode", flightCode).Return(mockSeats, nil)

	router := setupSeatRouter(mockService)

	httpRequest, _ := http.NewRequest("GET", fmt.Sprintf("/bookings/seats/%s", flightCode), nil)
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var seats []models.Seat
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &seats)
	assert.NoError(t, err)
	assert.Equal(t, seats, mockSeats)
	mockService.AssertExpectations(t)
}
