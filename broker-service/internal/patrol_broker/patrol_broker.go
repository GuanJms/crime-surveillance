package patrolbroker

import "brokerServiceApp/internal/patrol_broker/proto/patrolpb"

var pb *PatrolBroker

type PatrolBroker struct{}

func NewPatrolBroker() *PatrolBroker {
	if pb == nil {
		pb = &PatrolBroker{}
	}
	return pb
}

func PatrolsFromProto(patrolsProto []*patrolpb.Patrol) []*PatrolDTO {
	patrolsDTO := make([]*PatrolDTO, 0)

	for _, p := range patrolsProto {
		dto := &PatrolDTO{}
		dto.FromProto(p)
		patrolsDTO = append(patrolsDTO, dto)
	}
	return patrolsDTO
}
