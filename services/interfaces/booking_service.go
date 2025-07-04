package interfaces

import (
	"flyhorizons-bookingservice/models"
)

type BookingService interface {
	BookingExists(bookingID int) bool
	GetByUserID(userID int) []models.Booking
	Create(booking models.Booking) (*models.Booking, error)
	DeleteByBookingID(id int) (bool, error)
	Update(booking models.Booking) (*models.Booking, error)
}
