package main

import (
	"authServiceApp/data"
	"authServiceApp/handler"
	"authServiceApp/proto/authpb"
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

	authpb.RegisterAuthServiceServer(s, &handler.AuthServer{
		AuthModel: &data.AuthModel{
			Repo: app.Repo,
		},
	})

	log.Printf("gRPC Server started on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to listen for gRPC %v", err)
	}
}
