package main

import (
	"fmt"
	"log"
	"net"
	"patrolServiceApp/data"
	"patrolServiceApp/handler"
	"patrolServiceApp/proto/patrolpb"

	"google.golang.org/grpc"
)

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatal("Failed to listen for gRPC %v", err)
	}

	s := grpc.NewServer()

	patrolpb.RegisterPatrolServiceServer(s, &handler.PatrolServer{
		PatrolModel: &data.PatrolModel{
			Repo: app.Repo,
		},
	})

	log.Printf("gRPC Server started on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to listen for gRPC %v", err)
	}
}
