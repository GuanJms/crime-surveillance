package handler

import (
	"authServiceApp/data"
	"authServiceApp/proto/authpb"
	"authServiceApp/token"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	AuthModel *data.AuthModel
}

func (s *AuthServer) CreateNewUser(ctx context.Context, req *authpb.NewUserRequest) (*authpb.NewUserResponse, error) {
	// parse new user request into auth models
	newUser := data.NewUserFromRequest(req)

	// create new user in DB
	err := s.AuthModel.Repo.CreateUser(newUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create new user: %v", err)
	}
	token, err := token.GenerateToken(*newUser.ID, []string{*newUser.Role})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	// return token
	resp := &authpb.NewUserResponse{
		Token:   token,
		Success: true,
		Role:    "CITIZEN", // hardcoded base user
		Message: "Successfully created new user",
	}
	return resp, nil
}

func (s *AuthServer) UserLogin(ctx context.Context, cred *authpb.UserLoginCredentials) (*authpb.UserLoginResposne, error) {
	credentials := data.CredentialsFromRequest(cred)

	success, err := s.AuthModel.Repo.AuthenticateUser(credentials)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to authenticate user: %v", err)
	}

	if !success {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	user, err := s.AuthModel.Repo.GetUserInfo(credentials.Username)
	if err != nil {
		return nil, err
	}

	err = s.AuthModel.Repo.UpdateLoginTime(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update login time: %v", err)
	}

	token, err := token.GenerateToken(*user.ID, []string{*user.Role})

	if err != nil {
		return nil, err
	}

	return &authpb.UserLoginResposne{
		Token:   token,
		Role:    *user.Role,
		Success: true,
		Message: "successfully created a user login response",
	}, nil
}

func (s *AuthServer) ChangeUserRole(ctx context.Context, req *authpb.ChangeUserRoleRequest) (*authpb.ChangeUserRoleResponse, error) {
	id := req.Id
	role := req.Role.String()
	err := s.AuthModel.Repo.ChangeUserRoleTo(id, role)
	if err != nil {
		return nil, err
	}
	// successfully updated the role
	return &authpb.ChangeUserRoleResponse{
		Success: true,
		Message: "successfully updated the role",
	}, nil
}

func (s *AuthServer) GetAllUsers(ctx context.Context, req *authpb.GetAllUsersRequest) (*authpb.GetAllUsersResponse, error) {
	users, err := s.AuthModel.Repo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	usersProto := data.UsersToProto(users)

	return &authpb.GetAllUsersResponse{
		Users: usersProto,
	}, nil
}
