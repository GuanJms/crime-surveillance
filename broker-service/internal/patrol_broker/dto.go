package patrolbroker

import (
	"brokerServiceApp/internal/patrol_broker/proto/patrolpb"
	"brokerServiceApp/internal/ptr"
	"fmt"
	"log"
)

const nullLocation = -999.0

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

type UpdatePatrolStatusRequestDTO struct {
	PatrolId string
	Status   string `json:"status"`
}

type UpdatePatrolInfoRequestDTO struct {
	UserId      string   `json:"userID"`
	OfficerId   *string  `json:"officerId,omitempty"`
	OfficerName *string  `json:"officerName,omitempty"`
	Status      *string  `json:"status,omitempty"`
	Street      *string  `json:"street,omitempty"`
	City        *string  `json:"city,omitempty"`
	State       *string  `json:"state,omitempty"`
	Latitude    *float64 `json:"latitude,omitempty"`
	Longitude   *float64 `json:"longitude,omitempty"`
}

type UpdatePatrolLocationRequestDTO struct {
	UserID   string       `json:"userID"`
	Location *LocationDTO `json:"location"`
}

type AssignPatrolToCrimeRequestDTO struct {
	PatrolID string
	CrimeID  string
}

func (dto *PatrolDTO) ToProto() *patrolpb.Patrol {
	var status patrolpb.PatrolStatus
	var location *patrolpb.Location
	status = patrolpb.PatrolStatus(patrolpb.PatrolStatus_value[*dto.Status])
	if dto.Location != nil {
		location = &patrolpb.Location{
			Street: ptr.DeferOrZero(dto.Location.Street),
			City:   ptr.DeferOrZero(dto.Location.City),
			State:  ptr.DeferOrZero(dto.Location.State),
		}
	}

	return &patrolpb.Patrol{
		UserId:      dto.UserId,
		OfficerId:   ptr.DeferOrZero(dto.OfficerId),
		OfficerName: ptr.DeferOrZero(dto.OfficerName),
		Status:      &status,
		Location:    location,
	}
}

func (dto *PatrolDTO) FromProto(proto *patrolpb.Patrol) *PatrolDTO {
	if proto == nil {
		return nil
	}
	dto.UserId = proto.UserId
	dto.OfficerId = &proto.OfficerId
	dto.OfficerName = &proto.OfficerName

	if proto.Status != nil {
		status := proto.Status.String()
		dto.Status = &status
	}

	// log.Printf("SO FAR %v", p)

	if proto.Location != nil {
		dto.Location = &LocationDTO{
			Street: &proto.Location.Street,
			City:   &proto.Location.City,
			State:  &proto.Location.State,
		}
	}
	// log.Printf("SO FAR %v", p)
	return dto
}

func (dto *UpdatePatrolInfoRequestDTO) ToProto() *patrolpb.UpdatePatrolInfoRequest {
	if dto == nil {
		return nil
	}

	var status *patrolpb.PatrolStatus
	if dto.Status != nil {
		status = ptr.Of(patrolpb.PatrolStatus(patrolpb.PatrolStatus_value[*dto.Status]))
	}

	log.Printf("Status: %v", status)

	return &patrolpb.UpdatePatrolInfoRequest{
		UserId:      dto.UserId,
		OfficerId:   dto.OfficerId,
		OfficerName: dto.OfficerName,
		Status:      status,
		Street:      dto.Street,
		City:        dto.City,
		State:       dto.State,
		Latitude:    dto.Latitude,
		Longitude:   dto.Longitude,
	}
}

func (dto *UpdatePatrolStatusRequestDTO) ToProto() (*patrolpb.UpdatePatrolStatusRequest, error) {
	// log.Printf("Update patrol status request DTO - Received dto.Status: %v", dto.Status)
	value, ok := patrolpb.PatrolStatus_value[dto.Status]
	if !ok {
		return nil, fmt.Errorf("invalid status: %v", dto.Status)
	}
	status := patrolpb.PatrolStatus(value)
	// log.Printf("Converted %v and middle value: %v", status, patrolpb.PatrolStatus_value[dto.Status])
	return &patrolpb.UpdatePatrolStatusRequest{
		PatrolId: dto.PatrolId,
		Status:   status,
	}, nil
}

func (dto *UpdatePatrolLocationRequestDTO) ToProto() (*patrolpb.UpdatePatrolLocationRequest, error) {
	if dto == nil {
		return nil, fmt.Errorf("update patrol location request DTO is nil")
	}

	if dto.Location == nil {
		return nil, fmt.Errorf("update patrol location request DTO location is nil")
	}

	var latitude, longitude float64
	if dto.Location.Latitude == nil {
		latitude = nullLocation
	} else {
		latitude = *dto.Location.Latitude
	}

	if dto.Location.Longitude == nil {
		longitude = nullLocation
	} else {
		longitude = *dto.Location.Longitude
	}

	location := &patrolpb.Location{
		Street:    ptr.DeferOrZero(dto.Location.Street),
		City:      ptr.DeferOrZero(dto.Location.City),
		State:     ptr.DeferOrZero(dto.Location.State),
		Latitude:  latitude,
		Longitude: longitude,
	}

	return &patrolpb.UpdatePatrolLocationRequest{
		UserId:   dto.UserID,
		Location: location,
	}, nil
}

func (dto *AssignPatrolToCrimeRequestDTO) ToProto() (*patrolpb.AssignPatrolToCrimeRequest, error) {
	if dto == nil {
		return nil, fmt.Errorf("AssignPatrolToCrimeRequestDTO is empty")
	}

	return &patrolpb.AssignPatrolToCrimeRequest{
		CrimeId:  dto.CrimeID,
		PatrolId: dto.PatrolID,
	}, nil
}
