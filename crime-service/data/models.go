package data

import (
	"crimeServiceApp/proto/crimepb"
	"database/sql"
	"errors"
	"log"
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
	PatrolID    string    `json:"patrol_id,omitempty"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Location    Location  `json:"location"`
	ReportedAt  time.Time `json:"reported_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CrimeUpdate struct {
	ID             string          `json:"id"`
	ReporterID     *string         `json:"reporter_id"`
	PatrolID       *string         `json:"patrol_id,omitempty"`
	Description    *string         `json:"description"`
	Status         *string         `json:"status"`
	LocationUpdate *LocationUpdate `json:"location"`
}

type Location struct {
	Street    string  `json:"street"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// LocationUpdate is not working since Partial location update is not avaialbe yet
type LocationUpdate struct {
	Street    *string  `json:"street"`
	City      *string  `json:"city"`
	State     *string  `json:"state"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
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

func NewCrimeFromProtoUpdateRequest(req *crimepb.UpdateCrimeReportRequest) (*Crime, error) {
	if req == nil {
		return nil, errors.New("Crime report request is empty")
	}

	description := DeferOrZero(req.Description)
	location := DeferOrZero(req.Location)
	patrol_id := DeferOrZero(req.PatrolId)
	return &Crime{
		ID:          req.Id,
		ReporterID:  req.ReporterId,
		Description: description,
		PatrolID:    patrol_id,
		Status:      req.Status.String(),
		Location: Location{
			Street:    location.Street,
			City:      location.City,
			State:     location.State,
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
		},
	}, nil
}

// TODO: implement function that map UpdateCrimeReportRequest into UpdateCrime
func NewCrimeUpdateFromProtoUpdateRequest(req *crimepb.UpdateCrimeReportRequest) *CrimeUpdate {
	var statusStr *string
	var location *LocationUpdate

	if req.Status != nil {
		s := req.Status.String()
		statusStr = &s
	}
	if req.Location != nil {
		location = &LocationUpdate{
			Street:    &req.Location.Street,
			City:      &req.Location.City,
			State:     &req.Location.State,
			Latitude:  &req.Location.Latitude,
			Longitude: &req.Location.Longitude,
		}
	}
	log.Println("Creating crime update")

	return &CrimeUpdate{
		ID:             req.Id,
		ReporterID:     &req.ReporterId,
		Status:         statusStr,
		PatrolID:       req.PatrolId,
		Description:    req.Description,
		LocationUpdate: location,
	}

}

// Insert a new crime
func (cm *CrimeModel) Insert(c *Crime) {
	cm.Repo.InsertNewCrime(c)
}
