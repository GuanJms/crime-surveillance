package data

import (
	"crimeServiceApp/proto/crimepb"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

func NewRepository(conn *sql.DB) Repository {
	if conn == nil {
		return nil
	} else {
		return NewPostgresRepository(conn)
	}
}

type CrimeModel struct {
	Repo Repository
}

type Crime struct {
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

func NewCrimeFromProto(req *crimepb.CrimeReportRequest) (*Crime, error) {
	if req == nil {
		return nil, errors.New("Crime report request is empty")
	}
	return &Crime{
		ID:          uuid.New().String(),
		ReporterID:  req.ReporterId,
		Description: req.Description,
		Status:      "NEW",
		Location: Location{
			Street:    req.Location.Street,
			City:      req.Location.City,
			State:     req.Location.State,
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		},
	}, nil
}

// Insert a new crime
func (cm *CrimeModel) Insert(c *Crime) {
	cm.Repo.InsertNewCrime(c)
}
