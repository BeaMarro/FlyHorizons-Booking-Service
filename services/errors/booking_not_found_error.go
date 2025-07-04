package errors

import "fmt"

type BookingNotFoundError struct {
	ID int
}

func (e *BookingNotFoundError) Error() string {
	return fmt.Sprintf("Booking with the ID %d was not found", e.ID)
}

func NewBookingNotFoundError(id int, errorCode int) *BookingNotFoundError {
	return &BookingNotFoundError{ID: id}
}
