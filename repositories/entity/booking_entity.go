package entities

import "time"

type BookingEntity struct {
	ID          int               `gorm:"column:ID;primaryKey"`
	UserID      int               `gorm:"column:UserID"`
	FlightCode  string            `gorm:"column:FlightCode"`
	FlightClass int               `gorm:"column:FlightClass"`
	CreatedAt   time.Time         `gorm:"column:CreatedAt"`
	Passengers  []PassengerEntity `gorm:"foreignKey:BookingID;references:ID"` // One-to-many relationship
	Seats       []SeatEntity      `gorm:"foreignKey:BookingID;references:ID"` // One-to-many relationship
	Luggage     string            `gorm:"column:Luggage;type:string"`         // JSON list of integers (string)
	Status      string            `gorm:"column:Status"`
}

// Override the default table name
func (BookingEntity) TableName() string {
	return "Booking"
}
