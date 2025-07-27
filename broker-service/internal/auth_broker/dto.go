package authbroker

import (
	"brokerServiceApp/internal/auth_broker/proto/authpb"
	"fmt"
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
