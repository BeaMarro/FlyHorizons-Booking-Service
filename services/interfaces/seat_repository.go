package interfaces

import (
	entities "flyhorizons-bookingservice/repositories/entity"
)

type SeatRepository interface {
	GetByFlightCode(flightCode string) ([]entities.SeatOptionEntity, error)
}
