package model

import "time"

type ProjectWork struct {
	Name       string    `json:"Name"`
	NameOfWork string    `json:"NameOfWork"`
	Date       time.Time `json:"Date"`
	Type       string    `json:"Type"`
}
