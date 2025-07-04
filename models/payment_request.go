package models

type PaymentRequest struct {
	BookingID int     `json:"booking_id"`
	Payment   Payment `json:"payment"`
}
