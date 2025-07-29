package patrolbroker

import (
	"brokerServiceApp/internal/patrol_broker/proto/patrolpb"
	"brokerServiceApp/utils"
)

type PatrolDTO struct {
	UserId      string       `json:"userID"`
	OfficerId   *string      `json:"officerID"`
	OfficerName *string      `json:"officerName"`
	Status      *string      `json:"status"`
	Location    *LocationDTO `json:"location"`
}

type LocationDTO struct {
	Street    *string  `json:"street"`
	City      *string  `json:"city"`
	State     *string  `json:"state"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

func (p *PatrolDTO) ToProto() *patrolpb.Patrol {
	var status patrolpb.PatrolStatus
	var location *patrolpb.Location
	status = patrolpb.PatrolStatus(patrolpb.PatrolStatus_value[*p.Status])
	if p.Location != nil {
		location = &patrolpb.Location{
			Street: utils.DeferOrZero(p.Location.Street),
			City:   utils.DeferOrZero(p.Location.City),
			State:  utils.DeferOrZero(p.Location.State),
		}
	}

	return &patrolpb.Patrol{
		UserId:      p.UserId,
		OfficerId:   utils.DeferOrZero(p.OfficerId),
		OfficerName: utils.DeferOrZero(p.OfficerName),
		Status:      &status,
		Location:    location,
	}
}

func (p *PatrolDTO) FromProto(proto *patrolpb.Patrol) *PatrolDTO {
	if proto == nil {
		return nil
	}
	p.UserId = proto.UserId
	p.OfficerId = &proto.OfficerId
	p.OfficerName = &proto.OfficerName

	if proto.Status != nil {
		status := proto.Status.String()
		p.Status = &status
	}

	// log.Printf("SO FAR %v", p)

	if proto.Location != nil {
		p.Location = &LocationDTO{
			Street: &proto.Location.Street,
			City:   &proto.Location.City,
			State:  &proto.Location.State,
		}
	}
	// log.Printf("SO FAR %v", p)
	return p
}
