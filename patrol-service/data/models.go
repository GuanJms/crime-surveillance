package data

import (
	"errors"
	"fmt"
	"log"
	"patrolServiceApp/proto/patrolpb"
	"patrolServiceApp/ptr"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PatrolModel struct {
	Repo     Repository
	FastRepo FastRepository
}

func NewRepository(pool *pgxpool.Pool) Repository {
	if pool == nil {
		return nil
	} else {
		return NewPostgresRepository(pool)
	}
}

func NewFastRepository() FastRepository {
	return NewRedisRepository()
}

type Location struct {
	Street    *string    `json:"street,omitempty" db:"street" redis:"street"`
	City      *string    `json:"city,omitempty" db:"city" redis:"city"`
	State     *string    `json:"state,omitempty" db:"state" redis:"state"`
	Latitude  *float64   `json:"latitude,omitempty" db:"latitude" redis:"latitude"`
	Longitude *float64   `json:"longitude,omitempty" db:"longitude" redis:"longitude"`
	Timestamp *time.Time `reids:"timestamp"`
}

type Patrol struct {
	UserID      *string    `json:"userId,omitempty" db:"user_id"`
	OfficerID   *string    `json:"officerId,omitempty" db:"officer_id"`
	OfficerName *string    `json:"officerName,omitempty" db:"officer_name"`
	Status      *string    `json:"status,omitempty" db:"status"`
	Location    *Location  `json:"location,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}

type GetPatrolInfoRequest struct {
	UserID      *string `json:"userId,omitempty" db:"user_id"`
	OfficerID   *string `json:"officerId,omitempty" db:"officer_id"`
	OfficerName *string `json:"officerName,omitempty" db:"officer_name"`
	Status      *string `json:"status,omitempty" db:"status"`
	Street      *string `json:"street,omitempty" db:"street"`
	City        *string `json:"city,omitempty" db:"city"`
	State       *string `json:"state,omitempty" db:"state"`
}

type GetPatrolInfoResponse struct {
	Patrols []*Patrol `json:"patrols,omitempty"`
	Success *bool     `json:"success,omitempty"`
	Message *string   `json:"message,omitempty"`
}

type UpdatePatrolInfoRequest struct {
	UserId      string   `json:"userID" db:"user_id"`
	OfficerId   *string  `json:"officerId,omitempty" db:"officer_id"`
	OfficerName *string  `json:"officerName,omitempty" db:"officer_name"`
	Status      *string  `json:"status,omitempty" db:"status"`
	Street      *string  `json:"street,omitempty" db:"street"`
	City        *string  `json:"city,omitempty" db:"city"`
	State       *string  `json:"state,omitempty" db:"state"`
	Latitude    *float64 `json:"latitude,omitempty" db:"latitude"`
	Longitude   *float64 `json:"longitude,omitempty" db:"longitude"`
}

func (l *Location) ToProto() *patrolpb.Location {
	if l == nil {
		return nil
	}

	return &patrolpb.Location{
		Street:    ptr.Deref(l.Street),
		City:      ptr.Deref(l.City),
		State:     ptr.Deref(l.State),
		Latitude:  ptr.Deref(l.Latitude),
		Longitude: ptr.Deref(l.Longitude),
	}
}

func (p *Patrol) ToProto() (*patrolpb.Patrol, error) {
	if p == nil {
		return nil, nil
	}

	log.Printf("Convert Patrol to Proto")

	var status *patrolpb.PatrolStatus
	if p.Status != nil {
		status = ptr.Of(patrolpb.PatrolStatus(patrolpb.PatrolStatus_value[*p.Status]))
	}

	location := p.Location.ToProto()

	var createdAt *timestamppb.Timestamp
	if p.CreatedAt != nil {
		createdAt = timestamppb.New(*p.CreatedAt)
	}

	var updatedAt *timestamppb.Timestamp
	if p.UpdatedAt != nil {
		updatedAt = timestamppb.New(*p.UpdatedAt)
	}

	return &patrolpb.Patrol{
		UserId:      ptr.Deref(p.UserID),
		OfficerId:   ptr.Deref(p.OfficerID),
		OfficerName: ptr.Deref(p.OfficerName),
		Status:      status,
		Location:    location,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func PatrolFromProto(p *patrolpb.Patrol) (*Patrol, error) {
	if p == nil {
		return nil, nil
	}

	log.Printf("PatrolFromProto Received: %v", p)
	// convert patrol proto to patrol in data model
	var status *string
	if p.Status != nil {
		status = ptr.Of(p.Status.String())
	}

	var createdAt *time.Time
	if p.CreatedAt != nil {
		createdAt = ptr.Of(p.CreatedAt.AsTime())
	}

	var updatedAt *time.Time
	if p.UpdatedAt != nil {
		updatedAt = ptr.Of(p.UpdatedAt.AsTime())
	}

	var location *Location
	if p.Location != nil {
		location = &Location{
			Street:    ptr.Of(p.Location.Street),
			City:      ptr.Of(p.Location.City),
			State:     ptr.Of(p.Location.State),
			Latitude:  ptr.Of(p.Location.Latitude),
			Longitude: ptr.Of(p.Location.Longitude),
		}
	}

	return &Patrol{
		UserID:      ptr.Of(p.UserId),
		OfficerID:   ptr.Of(p.OfficerId),
		OfficerName: ptr.Of(p.OfficerName),
		Status:      status,
		Location:    location,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func (p *Patrol) String() string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf("Patrol{UserID: %v, OfficerID: %v, OfficerName: %v, Status: %v, Location: %v, CreatedAt: %v, UpdatedAt: %v}",
		ptr.Deref(p.UserID), ptr.Deref(p.OfficerID), ptr.Deref(p.OfficerName), ptr.Deref(p.Status), p.Location.String(), ptr.Deref(p.CreatedAt), ptr.Deref(p.UpdatedAt))
}

func (l *Location) String() string {
	if l == nil {
		return "nil"
	}
	return fmt.Sprintf("Location{Street: %v, City: %v, State: %v, Latitude: %v, Longitude: %v}",
		ptr.Deref(l.Street), ptr.Deref(l.City), ptr.Deref(l.State), ptr.Deref(l.Latitude), ptr.Deref(l.Longitude))
}

func GetPatrolInfoRequestFromProto(req *patrolpb.GetPatrolInfoRequest) *GetPatrolInfoRequest {
	return &GetPatrolInfoRequest{}
}

func PatrolsToProto(patrols []*Patrol) ([]*patrolpb.Patrol, error) {
	if patrols == nil {
		return nil, nil
	}
	patrolsProto := make([]*patrolpb.Patrol, 0)

	for _, p := range patrols {
		if p != nil {
			proto, err := p.ToProto()
			log.Printf("Converting to proto: %v", proto)
			if err != nil {
				return nil, errors.New("failed to convert patrols list to proto list")
			}
			patrolsProto = append(patrolsProto, proto)
		}
	}
	return patrolsProto, nil
}

func UpdatePatrolInfoRequestFromProto(req *patrolpb.UpdatePatrolInfoRequest) *UpdatePatrolInfoRequest {
	var status *string
	if req.Status != nil {
		status = ptr.Of(req.Status.String())
	}
	return &UpdatePatrolInfoRequest{
		UserId:      req.UserId,
		OfficerId:   req.OfficerId,
		OfficerName: req.OfficerName,
		Status:      status,
		Street:      req.Street,
		City:        req.City,
		State:       req.State,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
	}
}

func UpdatePatrolLocationRequestFromProto(req *patrolpb.UpdatePatrolLocationRequest) (string, *Location, error) {
	if req == nil {
		return "", nil, fmt.Errorf("update patrol location request is nil")
	}
	patrolId := req.UserId
	lp := req.Location
	if lp == nil {
		return patrolId, nil, fmt.Errorf("update patrol location request location is nil")
	}
	return patrolId, &Location{
		Street:    &lp.Street,
		City:      &lp.City,
		State:     &lp.State,
		Latitude:  &lp.Latitude,
		Longitude: &lp.Longitude,
	}, nil
}
