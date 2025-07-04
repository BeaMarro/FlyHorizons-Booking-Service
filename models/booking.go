package models

import (
	"flyhorizons-bookingservice/models/enums"
)

type Booking struct {
	ID          int               `json:"id"`
	UserID      int               `json:"user_id"`
	FlightCode  string            `json:"flight_code"`
	FlightClass enums.FlightClass `json:"flight_class"`
	Luggage     []enums.Luggage   `json:"luggage"`
	Seats       []Seat            `json:"seats"`
	Passengers  []Passenger       `json:"passengers"`
	Payment     Payment           `json:"payment"`
	Status      enums.Status      `json:"status"`
}
