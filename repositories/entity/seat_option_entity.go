package entities

type SeatOptionEntity struct {
	ID     int    `gorm:"column:ID;primaryKey;autoIncrement"`
	Row    int    `gorm:"column:row;not null"`
	Column string `gorm:"column:seat_column;type:char(1);not null"`
	Status bool   `gorm:"column:status"`
}

// Override the default table name
func (SeatOptionEntity) TableName() string {
	return "SeatOption"
}
