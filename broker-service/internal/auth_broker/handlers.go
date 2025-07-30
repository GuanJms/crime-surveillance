package authbroker

import (
	"brokerServiceApp/internal/auth_broker/proto/authpb"
	"brokerServiceApp/internal/authmiddleware"
	"brokerServiceApp/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthBrokerHandler struct {
	GrpcClient authpb.AuthServiceClient
	GrpcConn   *grpc.ClientConn
}

func NewAuthBrokerHandler() (*AuthBrokerHandler, error) {
	conn, err := grpc.NewClient("auth-service:50000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := authpb.NewAuthServiceClient(conn)
	return &AuthBrokerHandler{
		GrpcClient: client,
		GrpcConn:   conn,
	}, nil
}

func (h *AuthBrokerHandler) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	if h.GrpcConn == nil {
		utils.ErrorJSON(w, errors.New("authentication broker handler has no grpc connection"), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	newUserRequest := authpb.NewUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&newUserRequest); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	resp, err := h.GrpcClient.CreateNewUser(ctx, &newUserRequest) // resp - NewUserResposne
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	if !resp.Success {
		utils.ErrorJSON(w, errors.New(resp.Message), http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *AuthBrokerHandler) ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if h.GrpcConn == nil {
		utils.ErrorJSON(w, errors.New("authentication broker handler has no grpc connection"), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var reqDTO ChangeUserRoleRequestDTO
	reqDTO.ID = id

	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	req, err := reqDTO.toProto()
	if err != nil {
		http.Error(w, "Invalid JSON proto", http.StatusBadRequest)
		return
	}

	resp, err := h.GrpcClient.ChangeUserRole(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !resp.Success {
		http.Error(w, resp.Message, http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *AuthBrokerHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	if h.GrpcConn == nil {
		http.Error(w, "Authentication broker handler has no grpc connection", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var req authpb.UserLoginCredentials

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	resp, err := h.GrpcClient.UserLogin(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !resp.Success {
		http.Error(w, resp.Message, http.StatusUnauthorized)
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *AuthBrokerHandler) AddTo(r chi.Router) {
	r.Route("/users", func(users chi.Router) {
		users.Post("/login", h.UserLogin)
		users.Group(func(protected chi.Router) {
			protected.Use(authmiddleware.JWTMiddleware(utils.Secret))
			protected.With(authmiddleware.RequireRole("ADMIN")).Post("/register", h.CreateNewUser)
			// users.Post("/register", h.CreateNewUser)
		})

	})

	r.Group(func(protected chi.Router) {
		protected.Use(authmiddleware.JWTMiddleware(utils.Secret))
		protected.With(authmiddleware.RequireRole("ADMIN")).Patch("/admin/users/{id}/role", h.ChangeUserRole)
	})
}
