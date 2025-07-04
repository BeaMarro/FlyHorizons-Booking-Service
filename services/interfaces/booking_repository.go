package interfaces

import (
	"flyhorizons-bookingservice/models/enums"
	entities "flyhorizons-bookingservice/repositories/entity"
)

type BookingRepository interface {
	GetAll() []entities.BookingEntity
	GetByID(id int) entities.BookingEntity
	GetByUserID(userID int) []entities.BookingEntity
	Create(booking entities.BookingEntity) *entities.BookingEntity
	DeleteByBookingID(bookingID int) bool
	UpdateStatus(bookingID int, status enums.Status)
	Update(booking entities.BookingEntity) entities.BookingEntity
}
