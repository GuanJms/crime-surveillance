package patrolbroker

import (
	"brokerServiceApp/internal/authmiddleware"
	"brokerServiceApp/internal/patrol_broker/proto/patrolpb"
	"brokerServiceApp/utils"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PatrolBrokerHandler struct {
	GrpcClient patrolpb.PatrolServiceClient
	GrpcConn   *grpc.ClientConn
}

func NewPatrolBrokerHandler() (*PatrolBrokerHandler, error) {
	conn, err := grpc.NewClient("patrol-service:50002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := patrolpb.NewPatrolServiceClient(conn)
	return &PatrolBrokerHandler{
		GrpcClient: client,
		GrpcConn:   conn,
	}, nil
}

func (h *PatrolBrokerHandler) AddTo(r chi.Router) {
	r.Route("/patrols", func(patrols chi.Router) {
		patrols.Group(func(protected chi.Router) {
			protected.Use(authmiddleware.JWTMiddleware(utils.Secret))
			protected.With(authmiddleware.RequireRole("ADMIN")).Post("/register", h.RegisterNewPatrol)
			protected.With(authmiddleware.RequireRole("DISPATCHER", "ADMIN")).Get("/register", h.GetAllPatrolInfo)
		})
	})
}

func (h *PatrolBrokerHandler) RegisterNewPatrol(w http.ResponseWriter, r *http.Request) {
	if h.GrpcConn == nil {
		http.Error(w, "no gprc connection", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var p PatrolDTO

	json.NewDecoder(r.Body).Decode(&p)
	req := p.ToProto()
	// log.Printf("RegisterNewPatrol Received: %v", p)
	// log.Printf("RegisterNewPatrol Converted to Proto: %v", req)

	resp, err := h.GrpcClient.RegisterNewPatrol(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !resp.Success {
		utils.ErrorJSON(w, errors.New(resp.Message), http.StatusBadRequest)
		return
	}

	// successful regsitered
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *PatrolBrokerHandler) GetAllPatrolInfo(w http.ResponseWriter, r *http.Request) {
	if h.GrpcConn == nil {
		http.Error(w, "no gprc connection", http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := patrolpb.GetPatrolInfoRequest{}

	resp, err := h.GrpcClient.GetAllPatrolInfo(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !resp.Success {
		utils.ErrorJSON(w, errors.New(utils.DeferOrZero(resp.Message)), http.StatusBadRequest)
		return
	}

	patrols := PatrolsFromProto(resp.Patrols)

	// successful regsitered
	utils.WriteJSON(w, http.StatusOK, patrols)
}
func (h *PatrolBrokerHandler) UpdatePatrolLocation(w http.ResponseWriter, r *http.Request) {
	log.Panic("not implemented")
}
func (h *PatrolBrokerHandler) AssignPatrolToCrime(w http.ResponseWriter, r *http.Request) {
	log.Panic("not implemented")
}
func (h *PatrolBrokerHandler) UpdatePatrolStatus(w http.ResponseWriter, r *http.Request) {
	log.Panic("not implemented")
}
