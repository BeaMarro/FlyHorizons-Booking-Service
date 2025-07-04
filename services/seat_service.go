package services

import (
	"flyhorizons-bookingservice/models"
	"flyhorizons-bookingservice/services/converter"
	"flyhorizons-bookingservice/services/interfaces"
)

type SeatService struct {
	seatRepo      interfaces.SeatRepository
	seatConverter converter.SeatConverter
}

func NewSeatService(repo interfaces.SeatRepository, seatConverter converter.SeatConverter) *SeatService {
	return &SeatService{
		seatRepo:      repo,
		seatConverter: seatConverter,
	}
}

func (seatService *SeatService) GetByFlightCode(flightCode string) ([]models.Seat, error) {
	seatOptionEntities, err := seatService.seatRepo.GetByFlightCode(flightCode)
	seats := seatService.seatConverter.ConvertSeatOptionEntitiesToSeats(seatOptionEntities)
	return seats, err
}
