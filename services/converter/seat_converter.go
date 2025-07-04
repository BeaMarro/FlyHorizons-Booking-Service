package converter

import (
	"flyhorizons-bookingservice/models"
	entities "flyhorizons-bookingservice/repositories/entity"
	"fmt"
)

type SeatConverter struct {
}

func (seatConverter *SeatConverter) ConvertSeatEntitiesToSeats(seatEntities []entities.SeatEntity) []models.Seat {
	var seats []models.Seat
	for _, entity := range seatEntities {
		seats = append(seats, models.Seat{
			Row:       entity.Row,
			Column:    entity.Column,
			Available: true,
		})
	}
	return seats
}

func (seatConverter *SeatConverter) ConvertSeatsToSeatEntities(seats []models.Seat, bookingID int) []entities.SeatEntity {
	var seatEntities []entities.SeatEntity
	for _, seat := range seats {
		seatEntities = append(seatEntities, entities.SeatEntity{
			BookingID: bookingID,
			Row:       seat.Row,
			Column:    seat.Column,
		})
	}
	return seatEntities
}

func (seatConverter *SeatConverter) ConvertSeatOptionEntitiesToSeats(seatOptionEntities []entities.SeatOptionEntity) []models.Seat {
	var seats []models.Seat
	for _, entity := range seatOptionEntities {
		fmt.Printf("Entity - Row: %d, Column: %s, Status: %v\n", entity.Row, entity.Column, entity.Status)
		seats = append(seats, models.Seat{
			Row:       entity.Row,
			Column:    entity.Column,
			Available: entity.Status,
		})
	}
	return seats
}
