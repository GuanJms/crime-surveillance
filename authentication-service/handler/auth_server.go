package handler

import (
	"authServiceApp/data"
	"authServiceApp/proto/authpb"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	AuthModel *data.AuthModel
}
