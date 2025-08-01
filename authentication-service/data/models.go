package data

import (
	"authServiceApp/proto/authpb"
	"database/sql"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	newUserRole = "CITIZEN"
)

func NewRepository(conn *sql.DB) Repository {
	if conn == nil {
		return nil
	} else {
		return NewPostgresRepository(conn)
	}
}

type AuthModel struct {
	Repo Repository
}

type User struct {
	ID           *string    `json:"id,omitempty"`
	Username     *string    `json:"username,omitempty"`
	Password     *string    `json:"password,omitempty"`
	PasswordHash *string    `json:"password_hash,omitempty"`
	Role         *string    `json:"role,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
	LastActive   *time.Time `json:"last_active,omitempty"`
}

func NewUserFromRequest(req *authpb.NewUserRequest) *User {
	role := newUserRole
	return &User{
		Username: &req.Username,
		Password: &req.Password,
		Role:     &role,
	}
}

type UserLoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CredentialsFromRequest(req *authpb.UserLoginCredentials) *UserLoginCredentials {
	return &UserLoginCredentials{
		Username: req.Username,
		Password: req.Password,
	}
}

func UsersToProto(users []*User) []*authpb.User {
	usersProto := make([]*authpb.User, len(users))
	for i, user := range users {
		usersProto[i] = &authpb.User{
			Id:           DeferOrZero(user.ID),
			Username:     DeferOrZero(user.Username),
			Role:         authpb.UserRole(authpb.UserRole_value[DeferOrZero(user.Role)]),
			CreatedAt:    timestamppb.New(DeferOrZero(user.CreatedAt)),
			UpdatedAt:    timestamppb.New(DeferOrZero(user.UpdatedAt)),
			LastLogin:    timestamppb.New(DeferOrZero(user.LastLogin)),
			LastActivity: timestamppb.New(DeferOrZero(user.LastActive)),
		}
	}
	return usersProto
}
