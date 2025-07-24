package main

import (
	"crimeServiceApp/data"
	"crimeServiceApp/handler"
	"crimeServiceApp/proto/crimepb"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatal("Failed to listen for gRPC %v", err)
	}

	s := grpc.NewServer()

	crimepb.RegisterCrimeServiceServer(s, &handler.CrimeServer{
		CrimeModels: &data.CrimeModel{
			Repo: app.Repo,
		},
	})

	log.Printf("gRPC Server started on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to listen for gRPC %v", err)
	}
}
