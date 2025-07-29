package handler

import (
	"context"
	"log"
	"patrolServiceApp/data"
	"patrolServiceApp/proto/patrolpb"
	"patrolServiceApp/ptr"
)

type PatrolServer struct {
	patrolpb.UnimplementedPatrolServiceServer
	PatrolModel *data.PatrolModel
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
	patrolsProto, err := data.PatrolsToProto(patrols)
	if err != nil {
		return nil, err
	}
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
	panic("not implemented")
}

func (s *PatrolServer) AssignPatrolToCrime(ctx context.Context, req *patrolpb.AssignPatrolToCrimeRequest) (*patrolpb.AssignPatrolToCrimeResponse, error) {
	panic("not implemented")
}

func (s *PatrolServer) UpdatePatrolStatus(ctx context.Context, req *patrolpb.UpdatePatrolStatusRequest) (*patrolpb.UpdatePatrolStatusResponse, error) {
	panic("not implemented")
}
