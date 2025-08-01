package authbroker

import (
	"brokerServiceApp/internal/auth_broker/proto/authpb"
	"brokerServiceApp/internal/ptr"
	"fmt"
	"time"
)

type ChangeUserRoleRequestDTO struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

func (r *ChangeUserRoleRequestDTO) toProto() (*authpb.ChangeUserRoleRequest, error) {
	role, ok := authpb.UserRole_value[r.Role]
	if !ok {
		return nil, fmt.Errorf("invalid role: %s", r.Role)
	}
	return &authpb.ChangeUserRoleRequest{
		Id:   r.ID,
		Role: authpb.UserRole(role),
	}, nil
}

type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Role         string     `json:"role"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
	LastLogin    *time.Time `json:"lastLogin"`
	LastActivity *time.Time `json:"lastActivity"`
}

func UsersFromProto(users []*authpb.User) []*User {
	usersDTO := make([]*User, len(users))
	for i, user := range users {
		usersDTO[i] = &User{
			ID:           user.Id,
			Username:     user.Username,
			Role:         user.Role.String(),
			CreatedAt:    ptr.Of(user.CreatedAt.AsTime()),
			UpdatedAt:    ptr.Of(user.UpdatedAt.AsTime()),
			LastLogin:    ptr.Of(user.LastLogin.AsTime()),
			LastActivity: ptr.Of(user.LastActivity.AsTime()),
		}
	}
	return usersDTO
}
