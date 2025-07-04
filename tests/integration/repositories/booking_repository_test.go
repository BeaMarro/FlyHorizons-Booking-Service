package repositories_test

import (
	"flyhorizons-bookingservice/repositories"
	entities "flyhorizons-bookingservice/repositories/entity"
	"log"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestBookingRepository provides an in-memory SQLite database for testing
type TestBookingRepository struct {
	repositories.BaseRepository
}

func (repo *TestBookingRepository) CreateConnection() (*gorm.DB, error) {
	if repo.DB != nil {
		return repo.DB, nil
	}

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{}) // No shared cache
	if err != nil {
		log.Printf("Failed to create SQLite database: %v", err)
		return nil, err
	}

	// Enable foreign key support
	db.Exec("PRAGMA foreign_keys = ON")

	if err := db.AutoMigrate(&entities.BookingEntity{}, &entities.PassengerEntity{}, &entities.SeatEntity{}); err != nil {
		log.Printf("Failed to auto-migrate schema: %v", err)
		return nil, err
	}

	repo.DB = db
	return db, nil
}

func NewTestBookingRepository() *repositories.BookingRepository {
	baseRepo := &TestBookingRepository{}
	_, err := baseRepo.CreateConnection()
	if err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}
	return repositories.NewBookingRepository(&baseRepo.BaseRepository)
}

func getPassengerEntities() []entities.PassengerEntity {
	return []entities.PassengerEntity{
		{
			FullName:       "John Doe",
			DateOfBirth:    time.Date(1985, 7, 9, 1, 0, 0, 0, time.UTC),
			PassportNumber: "1234",
		},
		{
			FullName:       "Jane Doe",
			DateOfBirth:    time.Date(1986, 8, 8, 2, 30, 0, 0, time.UTC),
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

func getBookings(repo *repositories.BookingRepository) []entities.BookingEntity {
	testBookings := []entities.BookingEntity{
		{
			UserID:      2,
			FlightCode:  "FR788",
			FlightClass: 1,
			CreatedAt:   getDate(),
			Passengers:  getPassengerEntities(),
			Seats:       getSeatEntities(),
			Luggage:     getLuggageString(),
		},
		{
			UserID:      4,
			FlightCode:  "FR789",
			FlightClass: 2,
			CreatedAt:   getDate(),
			Passengers:  getPassengerEntities(),
			Seats:       getSeatEntities(),
			Luggage:     getLuggageString(),
		},
	}

	// Clear any existing data
	repo.DB.Exec("DELETE FROM Seat")
	repo.DB.Exec("DELETE FROM Passenger")
	repo.DB.Exec("DELETE FROM Booking")

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

func TestBookingRepositoryGetAllReturnsBookings(t *testing.T) {
	// Arrange
	bookingRepo := NewTestBookingRepository()
	testBookings := getBookings(bookingRepo)

	// Act
	bookings := bookingRepo.GetAll()

	// Assert
	assert.Equal(t, testBookings, bookings)
}

func TestBookingRepositoryGetByValidUserIDReturnsBookings(t *testing.T) {
	// Arrange
	bookingRepo := NewTestBookingRepository()
	userBookings := []entities.BookingEntity{getBookings(bookingRepo)[0]}
	userID := 2

	// Act
	bookings := bookingRepo.GetByUserID(userID)

	// Assert
	assert.Equal(t, userBookings, bookings)
}

func TestBookingRepositoryGetByInvalidUserIDReturnsNoBookings(t *testing.T) {
	// Arrange
	bookingRepo := NewTestBookingRepository()
	invalidUserID := 999

	// Act
	bookings := bookingRepo.GetByUserID(invalidUserID)

	// Assert
	assert.Equal(t, []entities.BookingEntity{}, bookings)
}

func TestBookingRepositoryCreateBookingReturnsNewBooking(t *testing.T) {
	// Arrange
	bookingRepo := NewTestBookingRepository()
	testBookings := getBookings(bookingRepo)
	bookingEntity := entities.BookingEntity{
		ID:          7,
		UserID:      1,
		FlightCode:  "FR787",
		FlightClass: 1,
		CreatedAt:   getDate(),
		Passengers:  getPassengerEntities(),
		Seats:       getSeatEntities(),
		Luggage:     getLuggageString(),
	}

	// Act
	booking := bookingRepo.Create(bookingEntity)
	bookings := bookingRepo.GetAll()

	// Assert
	assert.Len(t, bookings, len(testBookings)+1)
	assert.Equal(t, bookingEntity, *booking)
}

func TestBookingRepositoryDeleteByValidIDReturnsTrue(t *testing.T) {
	// Arrange
	bookingRepo := NewTestBookingRepository()
	testBookings := getBookings(bookingRepo)
	bookingID := testBookings[0].ID

	// Act
	isDeleted := bookingRepo.DeleteByBookingID(bookingID)
	bookings := bookingRepo.GetAll()

	// Assert
	assert.Len(t, bookings, len(testBookings)-1)
	assert.True(t, isDeleted)
}

func TestBookingRepositoryDeleteByInvalidIDReturnsFalse(t *testing.T) {
	// Arrange
	bookingRepo := NewTestBookingRepository()
	testBookings := getBookings(bookingRepo)
	invalidFlightCode := 999

	// Act
	isDeleted := bookingRepo.DeleteByBookingID(invalidFlightCode)
	bookings := bookingRepo.GetAll()

	// Assert
	assert.Len(t, bookings, len(testBookings))
	assert.False(t, isDeleted)
}

func TestBookingRepositoryUpdateValidBookingReturnsUpdatedBooking(t *testing.T) {
	// Arrange
	bookingRepo := NewTestBookingRepository()
	testBookings := getBookings(bookingRepo)
	// Update all booking fields
	updatedBooking := entities.BookingEntity{
		ID:          7,
		UserID:      1,
		FlightCode:  "FR787",
		FlightClass: 0,
		CreatedAt:   getDate(),
		Passengers:  []entities.PassengerEntity{getPassengerEntities()[0]},
		Seats:       []entities.SeatEntity{getSeatEntities()[0]},
		Luggage:     getLuggageString(),
	}

	// Act
	booking := bookingRepo.Update(updatedBooking)

	// Assert
	assert.Equal(t, updatedBooking, booking)
	assert.NotNil(t, testBookings)
}
