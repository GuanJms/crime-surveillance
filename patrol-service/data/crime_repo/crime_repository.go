package crimerepo

import (
	"context"
	"fmt"
	"log"
	"patrolServiceApp/data/crime_repo/proto/crimepb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CrimeRepository struct {
	grpcConn   *grpc.ClientConn
	grpcClient crimepb.CrimeServiceClient
}

func NewCrimeRepository(crimeService string) (*CrimeRepository, error) {
	conn, err := grpc.NewClient(crimeService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := crimepb.NewCrimeServiceClient(conn)
	return &CrimeRepository{
		grpcConn:   conn,
		grpcClient: client,
	}, nil
}

func (repo *CrimeRepository) UpdateCrimePatrolIDStatus(crimeId string, patrolId string, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	log.Printf("Sending update crime patrol id status request to crime service: %s, %s, %s", crimeId, patrolId, status)

	statusEnum := crimepb.CrimeStatus(crimepb.CrimeStatus_value[status])
	req := crimepb.UpdateCrimeReportRequest{
		Id:       crimeId,
		PatrolId: &patrolId,
		Status:   &statusEnum,
	}
	resp, err := repo.grpcClient.PatchCrime(ctx, &req)
	log.Printf("Update crime patrol id status response: %v", resp)
	if err != nil {
		return err
	} else if !resp.Successful {
		return fmt.Errorf("err: %v", resp.Message)
	}
	return nil
}
