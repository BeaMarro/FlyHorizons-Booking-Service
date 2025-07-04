package entities

import "time"

type PassengerEntity struct {
	ID             int           `gorm:"column:ID;primaryKey"`
	BookingID      int           `gorm:"column:BookingID;index"`             // Foreign key for the Booking table
	Booking        BookingEntity `gorm:"foreignKey:BookingID;references:ID"` // Relationship to BookingEntity
	FullName       string        `gorm:"column:FullName"`
	DateOfBirth    time.Time     `gorm:"column:DateOfBirth"`
	PassportNumber string        `gorm:"column:PassportNumber"`
	Email          string        `gorm:"column:Email"`
}

// Override the default table name
func (PassengerEntity) TableName() string {
	return "Passenger"
}
