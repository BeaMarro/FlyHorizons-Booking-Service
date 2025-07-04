package repositories

import (
	"flyhorizons-bookingservice/models/enums"
	entities "flyhorizons-bookingservice/repositories/entity"
	"flyhorizons-bookingservice/services/interfaces"
	"fmt"
	"log"
)

type BookingRepository struct {
	*BaseRepository
}

var _ interfaces.BookingRepository = (*BookingRepository)(nil)

func NewBookingRepository(baseRepo *BaseRepository) *BookingRepository {
	if baseRepo == nil {
		baseRepo = &BaseRepository{}
	}

	// Ensure the database connection is established
	_, err := baseRepo.CreateConnection()
	if err != nil {
		fmt.Printf("Failed to establish database connection: %v\n", err)
	} else {
		fmt.Println("Database connection successfully established")
	}

	return &BookingRepository{
		BaseRepository: baseRepo,
	}
}

func (repo *BookingRepository) GetAll() []entities.BookingEntity {
	var bookings []entities.BookingEntity

	repo.DB.Preload("Passengers").Preload("Seats").Find(&bookings)

	return bookings
}

func (repo *BookingRepository) GetByID(id int) entities.BookingEntity {
	db, _ := repo.CreateConnection()

	var booking entities.BookingEntity

	// This preloads the related Passengers and Seats
	db.Preload("Passengers").Preload("Seats").Where("ID = ?", id).Find(&booking)

	return booking
}

func (repo *BookingRepository) GetByUserID(userID int) []entities.BookingEntity {
	db, _ := repo.CreateConnection()

	var bookings []entities.BookingEntity

	// This preloads the related Passengers and Seats
	db.Preload("Passengers").Preload("Seats").Where("UserID = ?", userID).Find(&bookings)

	return bookings
}

func (repo *BookingRepository) Create(bookingEntity entities.BookingEntity) *entities.BookingEntity {
	db, _ := repo.CreateConnection()

	if err := db.Create(&bookingEntity).Error; err != nil {
		return nil
	}

	return &bookingEntity
}

func (repo *BookingRepository) DeleteByBookingID(bookingID int) bool {
	db, _ := repo.CreateConnection()

	// Delete associated passengers
	if err := db.Where("BookingID = ?", bookingID).Delete(&entities.PassengerEntity{}).Error; err != nil {
		log.Printf("Error deleting associated passengers: %v", err)
		return false
	}

	// Delete associated seats
	if err := db.Where("BookingID = ?", bookingID).Delete(&entities.SeatEntity{}).Error; err != nil {
		log.Printf("Error deleting associated seats: %v", err)
		return false
	}

	// Delete booking
	result := repo.DB.Delete(&entities.BookingEntity{}, bookingID)

	return result.Error == nil && result.RowsAffected > 0
}

func (repo *BookingRepository) UpdateStatus(bookingID int, status enums.Status) {
	db, _ := repo.CreateConnection()

	var bookingEntity entities.BookingEntity

	// Find booking by ID
	if err := db.First(&bookingEntity, bookingID).Error; err != nil {
		log.Printf("Booking not found with ID %d: %v", bookingID, err)
	}

	// Update status
	bookingEntity.Status = string(status)

	// Save updated entity
	if err := db.Save(&bookingEntity).Error; err != nil {
		log.Printf("Failed to update booking status: %v", err)
	}
}

func (repo *BookingRepository) Update(bookingEntity entities.BookingEntity) entities.BookingEntity {
	db, _ := repo.CreateConnection()

	db.Save(&bookingEntity)

	return bookingEntity
}
