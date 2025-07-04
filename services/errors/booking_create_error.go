package errors

import "fmt"

type BookingCreateError struct {
	ID int
}

func (e *BookingCreateError) Error() string {
	return fmt.Sprintf("Booking with the ID %d could not be created successfully", e.ID)
}

func NewBookingCreateError(id int, errorCode int) *BookingCreateError {
	return &BookingCreateError{ID: id}
}
