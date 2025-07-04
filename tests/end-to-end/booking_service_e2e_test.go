package endtoend

import (
	"bytes"
	"encoding/json"
	"flyhorizons-bookingservice/models"
	"flyhorizons-bookingservice/models/enums"
	"flyhorizons-bookingservice/repositories"
	entities "flyhorizons-bookingservice/repositories/entity"
	"flyhorizons-bookingservice/routes"
	"flyhorizons-bookingservice/services"
	"flyhorizons-bookingservice/services/converter"
	mock_repositories "flyhorizons-bookingservice/tests/mocks"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type BookingServiceEndToEndTests struct {
	repositories.BaseRepository
}

// Create a test version of BaseRepository that uses an in-memory SQLite database
func (repo *BookingServiceEndToEndTests) CreateConnection() (*gorm.DB, error) {
	if repo.DB != nil {
		return repo.DB, nil
	}

	// Use unique database name for each test instance to avoid conflicts
	dbName := "file::memory:?cache=shared&_txlock=immediate&_fk=1"
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to create SQLite database: %v", err)
		return nil, err
	}

	// Enable foreign key support and WAL mode for better concurrency
	db.Exec("PRAGMA foreign_keys = ON")
	db.Exec("PRAGMA journal_mode = WAL")

	if err := db.AutoMigrate(&entities.BookingEntity{}, &entities.PassengerEntity{}, &entities.SeatEntity{}); err != nil {
		log.Printf("Failed to auto-migrate schema: %v", err)
		return nil, err
	}

	repo.DB = db
	return db, nil
}

func NewTestBookingRepository() *repositories.BookingRepository {
	baseRepo := &BookingServiceEndToEndTests{}
	_, err := baseRepo.CreateConnection()
	if err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}
	return repositories.NewBookingRepository(&baseRepo.BaseRepository)
}

func cleanDatabase(repo *repositories.BookingRepository) {
	// Disable foreign key constraints temporarily for cleanup
	repo.DB.Exec("PRAGMA foreign_keys = OFF")

	// Delete in correct order to respect foreign key constraints
	repo.DB.Exec("DELETE FROM Seat")
	repo.DB.Exec("DELETE FROM Passenger")
	repo.DB.Exec("DELETE FROM Booking")

	// Reset auto-increment counters
	repo.DB.Exec("DELETE FROM sqlite_sequence WHERE name IN ('Seat', 'Passenger', 'Booking')")

	// Re-enable foreign key constraints
	repo.DB.Exec("PRAGMA foreign_keys = ON")
}

func getPassengerEntities() []entities.PassengerEntity {
	return []entities.PassengerEntity{
		{
			FullName:       "John Doe",
			DateOfBirth:    time.Date(1985, 7, 9, 1, 0, 0, 0, time.UTC),
			Email:          "john@doe.com",
			PassportNumber: "1234",
		},
		{
			FullName:       "Jane Doe",
			DateOfBirth:    time.Date(1986, 8, 8, 2, 30, 0, 0, time.UTC),
			Email:          "jane@doe.com",
			PassportNumber: "4321",
		},
	}
}

func getSeatEntities() []entities.SeatEntity {
	return []entities.SeatEntity{
		{Row: 1, Column: "A"},
		{Row: 1, Column: "B"},
	}
}

func getLuggageString() string {
	return `["SmallBag","Cargo20kg"]`
}

func getDate() time.Time {
	return time.Date(2025, 4, 3, 9, 0, 0, 0, time.UTC)
}

func setupBookings(repo *repositories.BookingRepository) []entities.BookingEntity {
	// Clean database first
	cleanDatabase(repo)

	testBookings := []entities.BookingEntity{
		{
			UserID:      1,
			FlightCode:  "FR788",
			FlightClass: 1,
			CreatedAt:   getDate(),
			Passengers:  getPassengerEntities(),
			Seats:       getSeatEntities(),
			Luggage:     getLuggageString(),
		},
		{
			UserID:      1,
			FlightCode:  "FR789",
			FlightClass: 2,
			CreatedAt:   getDate(),
			Passengers:  getPassengerEntities(),
			Seats:       getSeatEntities(),
			Luggage:     getLuggageString(),
		},
	}

	// Save bookings with associations
	for i := range testBookings {
		if err := repo.DB.Create(&testBookings[i]).Error; err != nil {
			log.Fatalf("Failed to create booking: %v", err)
		}

		if err := repo.DB.Preload("Passengers").Preload("Seats").First(&testBookings[i], testBookings[i].ID).Error; err != nil {
			log.Fatalf("Failed to fetch created booking: %v", err)
		}
	}

	return testBookings
}

// Setup
func setupBookingService(repo *repositories.BookingRepository) *services.BookingService {
	bookingConverter := converter.BookingConverter{}
	passengerConverter := converter.PassengerConverter{}
	seatConverter := converter.SeatConverter{}
	return services.NewBookingService(repo, bookingConverter, passengerConverter, seatConverter)
}

func setupBookingRouter(service services.BookingService, gatewayAuthMiddleware *mock_repositories.MockGatewayAuthMiddleware) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	routes.RegisterBookingRoutes(router, &service, gatewayAuthMiddleware)
	return router
}

func getFirstPassengers() []models.Passenger {
	return []models.Passenger{
		{
			ID:             1,
			FullName:       "John Doe",
			DateOfBirth:    time.Date(1985, 7, 9, 1, 0, 0, 0, time.UTC),
			Email:          "john@doe.com",
			PassportNumber: "1234",
		},
		{
			ID:             2,
			FullName:       "Jane Doe",
			DateOfBirth:    time.Date(1986, 8, 8, 2, 30, 0, 0, time.UTC),
			Email:          "jane@doe.com",
			PassportNumber: "4321",
		},
	}
}

func getSecondPassengers() []models.Passenger {
	return []models.Passenger{
		{
			ID:             3,
			FullName:       "John Doe",
			DateOfBirth:    time.Date(1985, 7, 9, 1, 0, 0, 0, time.UTC),
			Email:          "john@doe.com",
			PassportNumber: "1234",
		},
		{
			ID:             4,
			FullName:       "Jane Doe",
			DateOfBirth:    time.Date(1986, 8, 8, 2, 30, 0, 0, time.UTC),
			Email:          "jane@doe.com",
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

func getLuggages() []enums.Luggage {
	return []enums.Luggage{
		enums.SmallBag,
		enums.Cargo20kg,
	}
}

func getBookings() []models.Booking {
	return []models.Booking{
		{
			ID:          1,
			UserID:      1,
			FlightCode:  "FR788",
			FlightClass: 1,
			Passengers:  getFirstPassengers(),
			Seats:       getSeats(),
			Luggage:     getLuggages(),
		},
		{
			ID:          2,
			UserID:      1,
			FlightCode:  "FR789",
			FlightClass: 0,
			Passengers:  getSecondPassengers(),
			Seats:       getSeats(),
			Luggage:     getLuggages(),
		},
	}
}

func setupTestEnvironment(userID int) (*repositories.BookingRepository, *services.BookingService, *gin.Engine) {
	repo := NewTestBookingRepository()
	setupBookings(repo)
	service := setupBookingService(repo)
	mockMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", userID)
	router := setupBookingRouter(*service, mockMiddleware)
	return repo, service, router
}

// End-to-End Tests
func TestEndToEndGetAllBookingsByMatchingUserIDReturnsBookings(t *testing.T) {
	// Arrange
	_, _, router := setupTestEnvironment(1)
	mockBookings := getBookings()

	url := "/bookings/"
	httpRequest, _ := http.NewRequest("GET", url, nil)
	httpRequest.Header.Set("Authorization", "Bearer mocktoken12345")

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	// Unmarshal the JSON response
	var bookings []models.Booking
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &bookings)

	assert.NoError(t, err)
	assert.Equal(t, mockBookings, bookings)
}

func TestEndToEndCreateExistingBookingReturnsConflictError(t *testing.T) {
	// Arrange
	_, _, router := setupTestEnvironment(1)
	mockBooking := getBookings()[0]

	// Make the JSON to create the booking
	requestBody, _ := json.Marshal(mockBooking)
	httpRequest, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusConflict, responseRecorder.Code)
}

func TestEndToEndDeleteExistingBookingReturnsDeletedSuccessfully(t *testing.T) {
	// Arrange
	_, _, router := setupTestEnvironment(1)
	mockBooking := getBookings()[0]
	mockBookingID := mockBooking.ID

	url := fmt.Sprintf("/bookings/%v", mockBookingID)
	httpRequest, _ := http.NewRequest("DELETE", url, nil)
	httpRequest.Header.Set("Authorization", "Bearer mocktoken12345")

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestEndToEndDeleteNonExistingBookingReturnsNotFoundError(t *testing.T) {
	// Arrange
	_, _, router := setupTestEnvironment(1)
	mockBookingID := 999 // Use a clearly non-existent ID

	url := fmt.Sprintf("/bookings/%v", mockBookingID)
	httpRequest, _ := http.NewRequest("DELETE", url, nil)
	httpRequest.Header.Set("Authorization", "Bearer mocktoken12345")

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}

func TestEndToEndUpdateExistingBookingUsingMatchingUserIDReturnsUpdatedBooking(t *testing.T) {
	// Arrange
	_, _, router := setupTestEnvironment(1)
	mockBooking := getBookings()[0]

	// Make the JSON to update the booking
	requestBody, _ := json.Marshal(mockBooking)
	httpRequest, _ := http.NewRequest("PUT", "/bookings/", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	// Unmarshal the JSON response
	var booking models.Booking
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &booking)

	assert.NoError(t, err)
	assert.Equal(t, mockBooking, booking)
}

func TestEndToEndUpdateNonExistingBookingReturnsNotFoundError(t *testing.T) {
	// Arrange
	_, _, router := setupTestEnvironment(1)
	mockBooking := getBookings()[0]
	mockBooking.ID = 999 // Use a clearly non-existent ID

	// Make the JSON to update the booking
	requestBody, _ := json.Marshal(mockBooking)
	httpRequest, _ := http.NewRequest("PUT", "/bookings/", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}

func TestEndToEndUpdateBookingUsingNonMatchingUserIDReturnsAccessDenied(t *testing.T) {
	// Arrange
	_, _, router := setupTestEnvironment(3) // Different user ID
	mockBooking := getBookings()[0]

	// Make the JSON to update the booking
	requestBody, _ := json.Marshal(mockBooking)
	httpRequest, _ := http.NewRequest("PUT", "/bookings/", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
}
