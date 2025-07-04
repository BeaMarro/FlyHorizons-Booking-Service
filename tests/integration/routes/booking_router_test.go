package routes_test

import (
	"bytes"
	"encoding/json"
	"flyhorizons-bookingservice/models"
	"flyhorizons-bookingservice/models/enums"
	"flyhorizons-bookingservice/routes"
	"flyhorizons-bookingservice/services/errors"
	mock_repositories "flyhorizons-bookingservice/tests/mocks"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestBookingRoute struct {
}

// Setup
func setupBookingRouter(mockService *mock_repositories.MockBookingService, gatewayAuthMiddleware *mock_repositories.MockGatewayAuthMiddleware) *gin.Engine {
	router := gin.Default()

	routes.RegisterBookingRoutes(router, mockService, gatewayAuthMiddleware)

	return router
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

func getLuggageList() []enums.Luggage {
	return []enums.Luggage{enums.SmallBag, enums.Cargo20kg}
}

func getBookings() []models.Booking {
	return []models.Booking{
		{
			ID:          0,
			UserID:      2,
			FlightCode:  "FR788",
			FlightClass: 1,
			Luggage:     getLuggageList(),
			Seats:       getSeats(),
			Passengers:  getPassengers(),
		},
		{
			ID:          1,
			UserID:      4,
			FlightCode:  "FR789",
			FlightClass: 1,
			Luggage:     getLuggageList(),
			Seats:       getSeats(),
			Passengers:  getPassengers(),
		},
	}
}

// Router Integration Tests
func TestGetBookingsByUserIDMatchingWithBookingsReturnsBookingsJSON(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockBookingService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 4)
	bearerToken := "Bearer mocktoken12345"
	mockBookings := []models.Booking{getBookings()[1]}
	userID := 4
	mockService.On("GetByUserID", userID).Return(mockBookings)

	router := setupBookingRouter(mockService, mockAPIGatewayMiddleware)

	httpRequest, _ := http.NewRequest("GET", "/bookings/", nil)
	httpRequest.Header.Set("Authorization", bearerToken)
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var bookings []models.Booking
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &bookings)
	assert.NoError(t, err)
	assert.Equal(t, mockBookings, bookings)
	mockService.AssertExpectations(t)
}

func TestGetByMatchingUserIDWithoutBookingsReturnsValueNil(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockBookingService)
	mockUserID := 4
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", mockUserID)
	bearerToken := "Bearer mocktoken12345"
	mockService.On("GetByUserID", mockUserID).Return([]models.Booking{})

	router := setupBookingRouter(mockService, mockAPIGatewayMiddleware)

	httpRequest, _ := http.NewRequest("GET", "/bookings/", nil)
	httpRequest.Header.Set("Authorization", bearerToken)
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var bookings []models.Booking
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &bookings)
	assert.NoError(t, err)
	assert.Equal(t, []models.Booking{}, bookings)
	mockService.AssertExpectations(t)
}

func TestGetBookingsByUserIDNonMatchingWithBookingReturnsAccessDenied(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockBookingService)
	mockUserID := 999
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", mockUserID)
	bearerToken := "Bearer mocktoken12345"
	mockBookings := []models.Booking{getBookings()[1]}
	mockService.On("GetByUserID", mockUserID).Return(mockBookings)

	router := setupBookingRouter(mockService, mockAPIGatewayMiddleware)

	url := "/bookings/"
	httpRequest, _ := http.NewRequest("GET", url, nil)
	httpRequest.Header.Set("Authorization", bearerToken)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)

	var errResponse map[string]interface{}
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

// TODO: Fix
// func TestCreateNonExistingBookingUsingMatchingUserReturnsCreatedBooking(t *testing.T) {
// 	// Arrange
// 	mockService := new(mock_repositories.MockBookingService)
// 	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 4)
// 	bearerToken := "Bearer mocktoken12345"
// 	mockBooking := getBookings()[1]
// 	mockService.On("Create", mockBooking).Return(&mockBooking, nil)

// 	router := setupBookingRouter(mockService, mockAPIGatewayMiddleware)

// 	requestBody, _ := json.Marshal(mockBooking)
// 	httpRequest, _ := http.NewRequest("POST", "/bookings/", bytes.NewBuffer(requestBody))
// 	httpRequest.Header.Set("Content-Type", "application/json")
// 	httpRequest.Header.Set("Authorization", bearerToken)
// 	responseRecorder := httptest.NewRecorder()

// 	// Act
// 	router.ServeHTTP(responseRecorder, httpRequest)

// 	// Assert
// 	assert.Equal(t, http.StatusCreated, responseRecorder.Code)

// 	var booking models.Booking
// 	err := json.Unmarshal(responseRecorder.Body.Bytes(), &booking)
// 	assert.NoError(t, err)
// 	assert.Equal(t, mockBooking, booking)
// 	mockService.AssertExpectations(t)
// }

// TODO: Fix
// func TestCreateExistingBookingUsingMatchingUserReturnsHTTPStatusError(t *testing.T) {
// 	// Arrange
// 	mockService := new(mock_repositories.MockBookingService)
// 	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 4)
// 	bearerToken := "Bearer mocktoken12345"
// 	mockBooking := getBookings()[1]
// 	mockService.On("Create", mockBooking).Return(nil, errors.NewBookingExistsError(mockBooking.ID, 409))

// 	router := setupBookingRouter(mockService, mockAPIGatewayMiddleware)

// 	requestBody, _ := json.Marshal(mockBooking)

// 	httpRequest, _ := http.NewRequest("POST", "/bookings/", bytes.NewBuffer(requestBody))
// 	httpRequest.Header.Set("Content-Type", "application/json")
// 	httpRequest.Header.Set("Authorization", bearerToken)
// 	responseRecorder := httptest.NewRecorder()

// 	// Act
// 	router.ServeHTTP(responseRecorder, httpRequest)

// 	// Assert
// 	assert.Equal(t, http.StatusConflict, responseRecorder.Code)
// 	mockService.AssertExpectations(t)
// }

// TODO: Fix
// func TestCreateBookingUsingNonMatchingUserReturnsAccessDenied(t *testing.T) {
// 	// Arrange
// 	mockService := new(mock_repositories.MockBookingService)
// 	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 999)
// 	bearerToken := "Bearer mocktoken12345"
// 	mockBooking := getBookings()[1]

// 	router := setupBookingRouter(mockService, mockAPIGatewayMiddleware)

// 	requestBody, _ := json.Marshal(mockBooking)

// 	httpRequest, _ := http.NewRequest("POST", "/bookings/", bytes.NewBuffer(requestBody))
// 	httpRequest.Header.Set("Content-Type", "application/json")
// 	httpRequest.Header.Set("Authorization", bearerToken)
// 	responseRecorder := httptest.NewRecorder()

// 	// Act
// 	router.ServeHTTP(responseRecorder, httpRequest)

// 	// Assert
// 	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
// 	mockService.AssertExpectations(t)
// }

func TestDeleteExistingBookingReturnsHTTPStatusOK(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockBookingService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 4)
	bearerToken := "Bearer mocktoken12345"
	bookingID := getBookings()[0].ID
	mockService.On("DeleteByBookingID", bookingID).Return(true, nil)

	router := setupBookingRouter(mockService, mockAPIGatewayMiddleware)

	httpRequest, _ := http.NewRequest("DELETE", fmt.Sprintf("/bookings/%d", bookingID), nil)
	httpRequest.Header.Set("Authorization", bearerToken)
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestDeleteNonExistingBookingReturnsHTTPStatusError(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockBookingService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 4)
	bearerToken := "Bearer mocktoken12345"
	invalidBookingID := 999
	errorCode := 404
	mockService.On("DeleteByBookingID", invalidBookingID).Return(false, errors.NewBookingNotFoundError(invalidBookingID, errorCode))

	router := setupBookingRouter(mockService, mockAPIGatewayMiddleware)

	httpRequest, _ := http.NewRequest("DELETE", fmt.Sprintf("/bookings/%d", invalidBookingID), nil)
	httpRequest.Header.Set("Authorization", bearerToken)
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}

func TestUpdateExistingBookingUsingMatchingUserReturnsUpdatedUser(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockBookingService)
	userID := 2
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", userID)
	bearerToken := "Bearer mocktoken12345"
	mockBooking := getBookings()[0]
	mockBooking.FlightClass = 0
	mockBooking.Luggage = []enums.Luggage{
		enums.SmallBag,
		enums.CabinBag,
		enums.CabinBag,
		enums.Cargo20kg,
	}
	mockBooking.Seats = []models.Seat{
		{
			Row:       2,
			Column:    "E",
			Available: true,
		},
		{
			Row:       2,
			Column:    "F",
			Available: true,
		},
	}
	mockBooking.Passengers = []models.Passenger{
		{
			ID:             1,
			FullName:       "John G. Doe",
			DateOfBirth:    time.Date(1986, 7, 9, 1, 0, 0, 0, time.UTC),
			PassportNumber: "1234jjjj",
		},
		{
			ID:             2,
			FullName:       "Jane M. Doe",
			DateOfBirth:    time.Date(1985, 8, 8, 2, 30, 0, 0, time.UTC),
			PassportNumber: "4321gggg",
		},
	}
	mockService.On("Update", mockBooking).Return(&mockBooking, nil)

	router := setupBookingRouter(mockService, mockAPIGatewayMiddleware)

	requestBody, _ := json.Marshal(mockBooking)
	httpRequest, _ := http.NewRequest("PUT", "/bookings/", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", bearerToken)
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var booking models.Booking
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &booking)
	assert.NoError(t, err)
	assert.Equal(t, mockBooking, booking)
	mockService.AssertExpectations(t)
}

func TestUpdateNonExistingBookingUsingMatchingUserReturnsHTTPStatusError(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockBookingService)
	userID := 2
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", userID)
	bearerToken := "Bearer mocktoken12345"
	mockBooking := getBookings()[0]
	errorCode := 404
	mockService.On("Update", mockBooking).Return(nil, errors.NewBookingNotFoundError(mockBooking.ID, errorCode))

	router := setupBookingRouter(mockService, mockAPIGatewayMiddleware)

	requestBody, _ := json.Marshal(mockBooking)
	httpRequest, _ := http.NewRequest("PUT", "/bookings/", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", bearerToken)
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateBookingUsingNonMatchingUserReturnsAccessDenied(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockBookingService)
	userID := 999
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", userID)
	bearerToken := "Bearer mocktoken12345"
	mockBooking := getBookings()[0]

	router := setupBookingRouter(mockService, mockAPIGatewayMiddleware)

	requestBody, _ := json.Marshal(mockBooking)
	httpRequest, _ := http.NewRequest("PUT", "/bookings/", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", bearerToken)
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
}
