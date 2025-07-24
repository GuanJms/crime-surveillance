package crimebroker

import "time"

type CrimeDTO struct {
	ID          string    `json:"id"`
	ReporterID  string    `json:"reporter_id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Location    Location  `json:"location"`
	ReportedAt  time.Time `json:"reported_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Location struct {
	Street    string  `json:"street"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
