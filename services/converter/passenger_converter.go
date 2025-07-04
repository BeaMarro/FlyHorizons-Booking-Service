package converter

import (
	"flyhorizons-bookingservice/models"
	entities "flyhorizons-bookingservice/repositories/entity"
)

type PassengerConverter struct {
}

func (passengerConverter *PassengerConverter) ConvertPassengerEntitiesToPassengers(passengerEntities []entities.PassengerEntity) []models.Passenger {
	var passengers []models.Passenger
	for _, entity := range passengerEntities {
		passengers = append(passengers, models.Passenger{
			ID:             entity.ID,
			FullName:       entity.FullName,
			DateOfBirth:    entity.DateOfBirth,
			PassportNumber: entity.PassportNumber,
			Email:          entity.Email,
		})
	}
	return passengers
}

func (passengerConverter *PassengerConverter) ConvertPassengersToPassengerEntities(passengers []models.Passenger, bookingID int) []entities.PassengerEntity {
	var passengerEntities []entities.PassengerEntity
	for _, passenger := range passengers {
		passengerEntities = append(passengerEntities, entities.PassengerEntity{
			ID:             passenger.ID,
			BookingID:      bookingID,
			FullName:       passenger.FullName,
			DateOfBirth:    passenger.DateOfBirth,
			PassportNumber: passenger.PassportNumber,
			Email:          passenger.Email,
		})
	}
	return passengerEntities
}
