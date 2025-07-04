package models

type Seat struct {
	Row       int    `json:"row"`
	Column    string `json:"column"`
	Available bool   `json:"available"`
}
