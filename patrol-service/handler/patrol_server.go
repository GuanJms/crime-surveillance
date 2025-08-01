package handler

import (
	"context"
	"fmt"
	"log"
	"patrolServiceApp/data"
	"patrolServiceApp/proto/patrolpb"
	"patrolServiceApp/ptr"
	"time"
)

const (
	syncPatrolLocationInterval = time.Second * 10
)

type PatrolServer struct {
	patrolpb.UnimplementedPatrolServiceServer
	PatrolModel *data.PatrolModel
	CrimeModel  *data.CrimeModel
}

func (s *PatrolServer) Init() {
	s.startSyncPatrolLocation(syncPatrolLocationInterval)
}

func (s *PatrolServer) RegisterNewPatrol(ctx context.Context, req *patrolpb.Patrol) (*patrolpb.NewPatrolRegisterResponse, error) {
	patrol, err := data.PatrolFromProto(req)
	log.Printf("Register new patrol conveted to patrol: %v", patrol)
	if err != nil {
		return nil, err
	}
	id, err := s.PatrolModel.Repo.InsertPatrol(patrol)
	if err != nil {
		return nil, err
	}

	err = s.PatrolModel.Repo.UpdateUserRole(id, "PATROL")
	if err != nil {
		return nil, err
	}
	// successful
	var resp patrolpb.NewPatrolRegisterResponse
	resp.PatrolId = id
	resp.Success = true
	resp.Message = "successfully registered patrol"

	return &resp, nil
}

func (s *PatrolServer) GetAllPatrolInfo(ctx context.Context, req *patrolpb.GetPatrolInfoRequest) (*patrolpb.GetPatrolInfoResponse, error) {
	// TODO: Adding GetPatrolInfoRequest
	reqDTO := data.GetPatrolInfoRequestFromProto(req)
	patrols, err := s.PatrolModel.Repo.GetAllPatrols(reqDTO)
	if err != nil {
		return nil, err
	}
	// log.Printf("Patrols: %v", patrols)
	patrolsProto, err := data.PatrolsToProto(patrols)
	if err != nil {
		return nil, err
	}
	// log.Printf("PatrolsProto: %v", patrolsProto)
	return &patrolpb.GetPatrolInfoResponse{
		Patrols: patrolsProto,
		Success: true,
		Message: ptr.Of("successfully returned all patrol information"),
	}, nil
}

func (s *PatrolServer) PutPatrolInfo(ctx context.Context, req *patrolpb.UpdatePatrolInfoRequest) (*patrolpb.UpdatePatrolInfoResponse, error) {
	reqDTO := data.UpdatePatrolInfoRequestFromProto(req)
	err := s.PatrolModel.Repo.PutPatrolInfo(reqDTO)
	if err != nil {
		return nil, err
	}
	return &patrolpb.UpdatePatrolInfoResponse{
		PatrolId: req.UserId,
		Success:  true,
		Message:  ptr.Of("successfully updated the patrol information"),
	}, nil
}

func (s *PatrolServer) PatchPatrolInfo(ctx context.Context, req *patrolpb.UpdatePatrolInfoRequest) (*patrolpb.UpdatePatrolInfoResponse, error) {
	reqDTO := data.UpdatePatrolInfoRequestFromProto(req)
	err := s.PatrolModel.Repo.PatchPatrolInfo(reqDTO)
	if err != nil {
		return nil, err
	}
	return &patrolpb.UpdatePatrolInfoResponse{
		PatrolId: req.UserId,
		Success:  true,
		Message:  ptr.Of("successfully updated the patrol information"),
	}, nil
}

func (s *PatrolServer) UpdatePatrolLocation(ctx context.Context, req *patrolpb.UpdatePatrolLocationRequest) (*patrolpb.UpdatePatrolLocationResponse, error) {
	patrolId, location, err := data.UpdatePatrolLocationRequestFromProto(req)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, fmt.Errorf("no location available")
	}

	// log.Printf("Updating patrol location for patrolId: %s, location: %v", patrolId, location)
	err = s.PatrolModel.FastRepo.UpdatePatrolLocation(patrolId, location)
	if err != nil {
		return nil, err
	}

	return &patrolpb.UpdatePatrolLocationResponse{
		UserId:  patrolId,
		Success: true,
	}, nil

}

func (s *PatrolServer) AssignPatrolToCrime(ctx context.Context, req *patrolpb.AssignPatrolToCrimeRequest) (*patrolpb.AssignPatrolToCrimeResponse, error) {
	patrolId := req.PatrolId
	crimeId := req.CrimeId

	log.Printf("Assigning patrol to crime: %s, %s", patrolId, crimeId)

	// Update status of patrol, conditioning with status as available
	err := s.PatrolModel.Repo.UpdatePatrolStatus(patrolId, "BUSY", ptr.Of("AVAILABLE"))
	if err != nil {
		return nil, err
	}

	// Update crime status and patrolid in table crime
	err = s.CrimeModel.Repo.UpdateCrimePatrolIDStatus(crimeId, patrolId, "ASSIGNED")
	if err != nil {
		return nil, err
	}

	return &patrolpb.AssignPatrolToCrimeResponse{
		PatrolId: patrolId,
		CrimeId:  crimeId,
		Success:  true,
		Message:  ptr.Of("successfully assigned patrol to crime"),
	}, nil
}

func (s *PatrolServer) UpdatePatrolStatus(ctx context.Context, req *patrolpb.UpdatePatrolStatusRequest) (*patrolpb.UpdatePatrolStatusResponse, error) {
	err := s.PatrolModel.Repo.UpdatePatrolStatus(req.PatrolId, req.Status.String(), nil)
	if err != nil {
		return nil, err
	}
	return &patrolpb.UpdatePatrolStatusResponse{
		PatrolId: req.PatrolId,
		Success:  true,
		Message:  ptr.Of("successfully updated status"),
	}, nil
}

func (s *PatrolServer) startSyncPatrolLocation(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		log.Printf("Sync patrol location started")
		defer ticker.Stop() // keep ticker alive durign go sync timeline
		for range ticker.C {
			if err := s.PatrolModel.FastRepo.SyncPatrolLocation(); err != nil {
				log.Printf("Syn patrol location error: ", err)
			}
		}
		log.Printf("Sync patrol location stopped")
	}()
}
