package entities

type SeatEntity struct {
	ID        int           `gorm:"column:ID;primaryKey"`
	BookingID int           `gorm:"column:BookingID;index"`             // Foreign key for the Booking table
	Booking   BookingEntity `gorm:"foreignKey:BookingID;references:ID"` // Relationship to BookingEntity
	Row       int           `gorm:"column:Row"`
	Column    string        `gorm:"column:Column"`
}

// Override the default table name
func (SeatEntity) TableName() string {
	return "Seat"
}
