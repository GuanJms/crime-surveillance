package mapper

import (
	"crimeServiceApp/data"
	"crimeServiceApp/proto/crimepb"
	"errors"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func fromPostgresStatus(s string) (crimepb.CrimeStatus, error) {
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

func ToProtoCrime(c *data.Crime) (*crimepb.Crime, error) {
	// converting crime to proto crime
	status, err := fromPostgresStatus(c.Status)
	if err != nil {
		return nil, err
	}
	return &crimepb.Crime{
		Id:          c.ID,
		ReporterId:  c.ReporterID,
		Description: c.Description,
		Status:      status,
		Location: &crimepb.Location{
			Street:    c.Location.Street,
			City:      c.Location.City,
			State:     c.Location.State,
			Latitude:  c.Location.Latitude,
			Longitude: c.Location.Longitude,
		},
		ReportedAt: timestamppb.New(c.ReportedAt),
		CreatedAt:  timestamppb.New(c.CreatedAt),
		UpdatedAt:  timestamppb.New(c.UpdatedAt),
	}, nil
}

func FromProtoCrime(pc *crimepb.Crime) (*data.Crime, error) {
	panic("Not implemented")
}
