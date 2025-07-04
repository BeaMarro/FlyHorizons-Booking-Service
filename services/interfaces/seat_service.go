package interfaces

import (
	"flyhorizons-bookingservice/models"
)

type SeatService interface {
	GetByFlightCode(flightCode string) ([]models.Seat, error)
}
