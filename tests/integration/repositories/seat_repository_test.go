package repositories_test

import (
	"flyhorizons-bookingservice/repositories"
	entities "flyhorizons-bookingservice/repositories/entity"
	"log"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Create a test version of TestSeatRepository that uses an in-memory SQLite database
type TestSeatRepository struct {
	repositories.BaseRepository
}

func (repo *TestSeatRepository) CreateConnection() (*gorm.DB, error) {
	if repo.DB != nil {
		return repo.DB, nil
	}

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to create SQLite database: %v", err)
		return nil, err
	}

	// Auto migrate entities for the test database
	if err := db.AutoMigrate(&entities.SeatOptionEntity{}, &entities.SeatEntity{}); err != nil {
		log.Printf("Failed to auto-migrate schema: %v", err)
		return nil, err
	}

	repo.DB = db
	return db, nil
}

func NewTestSeatRepository() *repositories.SeatRepository {
	baseRepo := &TestBookingRepository{}
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{}) // No shared cache
	if err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}

	// Auto-migrate tables for the test database
	if err := db.AutoMigrate(&entities.SeatOptionEntity{}); err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	baseRepo.DB = db
	return repositories.NewSeatRepository(&baseRepo.BaseRepository)
}

func getSeatOptionEntities(repo *repositories.SeatRepository) []entities.SeatOptionEntity {
	// seats := []entities.SeatEntity{
	// 	{
	// 		ID:        1,
	// 		BookingID: getBookings(repositories.NewBookingRepository())[0].ID,
	// 		Booking:   getBookings(repositories.NewBookingRepository())[0],
	// 		Row:       1,
	// 		Column:    "A",
	// 	},
	// 	{
	// 		ID:        2,
	// 		BookingID: getBookings(repositories.NewBookingRepository())[1].ID,
	// 		Booking:   getBookings(repositories.NewBookingRepository())[1],
	// 		Row:       1,
	// 		Column:    "B",
	// 	},
	// }

	seatOptionEntities := []entities.SeatOptionEntity{
		{
			ID:     1,
			Row:    1,
			Column: "A",
			Status: true,
		},
		{
			ID:     1,
			Row:    1,
			Column: "B",
			Status: false,
		},
	}

	// Clear any existing data
	repo.DB.Exec("DELETE FROM Seat")

	// Save bookings with associations
	// Create Seats in the Seat Table
	// Create SeatOptions in the SeatOptions Table

	return seatOptionEntities
}

func TestSeatRepositoryGetByValidFlightCodeReturnsSeats(t *testing.T) {
	// Arrange
	// seatRepo := NewTestSeatRepository()
	// flightSeats := getSeatEntities()
	// flightCode := 1

	// // Act
	// seats, err := seatRepo.GetByFlightID(flightCode)

	// // Assert
	// assert.NoError(t, err)
	// assert.Equal(t, flightSeats, seats)
}

// Returns all seats
// But if none have been booked yet, all of them are available, per say for a whole new flight
func TestSeatRepositoryGetByInvalidFlightCodeReturnsNoSeats(t *testing.T) {
	// Arrange

	// Act

	// Assert
}
