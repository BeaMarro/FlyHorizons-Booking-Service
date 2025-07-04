package repositories

import (
	entities "flyhorizons-bookingservice/repositories/entity"
	"fmt"
)

type SeatRepository struct {
	*BaseRepository
}

func NewSeatRepository(baseRepo *BaseRepository) *SeatRepository {
	return &SeatRepository{
		BaseRepository: baseRepo,
	}
}

func (r *SeatRepository) GetByFlightCode(flightCode string) ([]entities.SeatOptionEntity, error) {
	db, err := r.CreateConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to create DB connection: %v", err)
	}

	var results []entities.SeatOptionEntity

	err = db.Raw(`
		SELECT 
			so.Row AS row,
			so.[Column] AS seat_column,
			CASE 
				WHEN s.BookingID IS NULL THEN 1  -- Available (no booking)
				ELSE 0  -- Booked (has a booking)
			END AS status
		FROM 
			SeatOption so
		LEFT JOIN 
			Seat s ON so.Row = s.Row AND so.[Column] = s.[Column]
		LEFT JOIN 
			Booking b ON s.BookingID = b.ID 
		WHERE 
			(b.FlightCode = ? OR b.FlightCode IS NULL)  -- Select only bookings for the FlightID
		ORDER BY 
			so.Row, so.[Column]
	`, flightCode).Scan(&results).Error

	if err != nil {
		fmt.Printf("Error fetching seat options for flightCode %s: %v\n", flightCode, err)
		return nil, err
	}

	return results, nil
}
