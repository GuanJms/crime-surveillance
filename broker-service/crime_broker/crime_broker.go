package crimebroker

import "brokerServiceApp/crime_broker/proto/crimepb"

var cb *CrimeBroker

type CrimeBroker struct{}

func NewCrimeBroker() *CrimeBroker {
	// if the crime broker is not initialized, initialize it
	if cb == nil {
		cb = &CrimeBroker{}
	}
	return cb
}

func CrimeProtoToDTO(c *crimepb.Crime) *CrimeDTO {
	return &CrimeDTO{
		ID:          c.Id,
		ReporterID:  c.ReporterId,
		Description: c.Description,
		Status:      c.Status.String(),
		Location: Location{
			Street:    c.Location.Street,
			City:      c.Location.City,
			State:     c.Location.State,
			Latitude:  c.Location.Latitude,
			Longitude: c.Location.Longitude,
		},
		ReportedAt: c.ReportedAt.AsTime(),
		CreatedAt:  c.CreatedAt.AsTime(),
		UpdatedAt:  c.UpdatedAt.AsTime(),
	}
}

func CrimesProtoToDTOs(list []*crimepb.Crime) []CrimeDTO {
	dtos := make([]CrimeDTO, 0, len(list))
	for _, c := range list {
		if c != nil {
			dtos = append(dtos, *CrimeProtoToDTO(c))
		}
	}
	return dtos
}
