package models

import "time"

type Passenger struct {
	ID             int       `json:"id"`
	FullName       string    `json:"full_name"`
	DateOfBirth    time.Time `json:"date_of_birth"`
	PassportNumber string    `json:"passport_number"`
	Email          string    `json:"email"`
}
