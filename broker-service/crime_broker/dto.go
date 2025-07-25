package crimebroker

import (
	"brokerServiceApp/crime_broker/proto/crimepb"
	"brokerServiceApp/utils"
	"errors"
	"strings"
	"time"
)

type CrimeDTO struct {
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

type Location struct {
	Street    string  `json:"street"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func fromPostgresStatusToGRPCCrimeStatus(s string) (crimepb.CrimeStatus, error) {
	switch strings.ToUpper(s) {
	case "NEW":
		return crimepb.CrimeStatus_NEW, nil
	case "ASSIGNED":
		return crimepb.CrimeStatus_ASSIGNED, nil
	case "RESOLVED":
		return crimepb.CrimeStatus_RESOLVED, nil
	default:
		return crimepb.CrimeStatus_NEW, errors.New("unsupported crime status")
	}
}

type UpdateCrimeReportRequestDTO struct {
	Id          *string   `json:"id"`
	ReporterId  *string   `json:"reporter_id,omitempty"`
	PatrolId    *string   `json:"patrol_id,omitempty"`
	Description *string   `json:"description"`
	Status      *string   `json:"status"`
	Location    *Location `json:"location"`
}

func (dto *UpdateCrimeReportRequestDTO) toProto() (*crimepb.UpdateCrimeReportRequest, error) {
	status, err := fromPostgresStatusToGRPCCrimeStatus(*dto.Status)
	if err != nil {
		return nil, err
	}
	var req crimepb.UpdateCrimeReportRequest
	req.Id = utils.DeferOrZero(dto.Id)
	req.ReporterId = utils.DeferOrZero(dto.ReporterId)
	req.PatrolId = dto.PatrolId
	req.Description = dto.Description
	req.Status = &status

	if dto.Location != nil {
		req.Location = &crimepb.Location{
			Street:    dto.Location.Street,
			City:      dto.Location.City,
			State:     dto.Location.State,
			Latitude:  dto.Location.Latitude,
			Longitude: dto.Location.Longitude,
		}
	}
	return &req, nil
}
