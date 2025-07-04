package errors

import "fmt"

type BookingExistsError struct {
	ID int
}

func (e *BookingExistsError) Error() string {
	return fmt.Sprintf("Booking with the ID %d already exists", e.ID)
}

func NewBookingExistsError(id int, errorCode int) *BookingExistsError {
	return &BookingExistsError{ID: id}
}
