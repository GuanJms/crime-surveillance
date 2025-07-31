package patrolbroker

import (
	"brokerServiceApp/internal/authmiddleware"
	"brokerServiceApp/internal/patrol_broker/proto/patrolpb"
	"brokerServiceApp/internal/ptr"
	"brokerServiceApp/internal/utils"
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
		patrols.Group(func(p chi.Router) {
			p.Use(authmiddleware.JWTMiddleware(utils.Secret))
			p.With(authmiddleware.RequireRole("ADMIN")).Post("/register", h.RegisterNewPatrol)
			p.With(authmiddleware.RequireRole("DISPATCHER", "ADMIN")).Get("/register", h.GetAllPatrolInfo)

			p.With(authmiddleware.RequireRole("PATROL")).Put("/", h.PutPatrolInfo)
			p.With(authmiddleware.RequireRole("PATROL")).Patch("/", h.PatchPatrolInfo)
			p.With(authmiddleware.RequireRole("PATROL")).Put("/status", h.UpdatePatrolStatus)
			p.With(authmiddleware.RequireRole("PATROL")).Put("/location", h.UpdatePatrolLocation)
			p.With(authmiddleware.RequireRole("DISPATCHER", "ADMIN")).Post("/dispatch", h.AssignPatrolToCrime)
		})

		patrols.Route("/{id}", func(idRouter chi.Router) {
			idRouter.Use(authmiddleware.JWTMiddleware(utils.Secret))
			idRouter.With(authmiddleware.RequireRole("ADMIN")).Put("/", h.PutPatrolInfo)
			idRouter.With(authmiddleware.RequireRole("ADMIN")).Patch("/", h.PatchPatrolInfo)
			idRouter.With(authmiddleware.RequireRole("DISPATCHER", "ADMIN")).Put("/status", h.UpdatePatrolStatus)
			idRouter.With(authmiddleware.RequireRole("DISPATCHER", "ADMIN")).Put("/location", h.UpdatePatrolLocation)
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
		utils.ErrorJSON(w, errors.New(ptr.DeferOrZero(resp.Message)), http.StatusBadRequest)
		return
	}

	patrols := PatrolsFromProto(resp.Patrols)

	// successful regsitered
	utils.WriteJSON(w, http.StatusOK, patrols)
}

func (h *PatrolBrokerHandler) PutPatrolInfo(w http.ResponseWriter, r *http.Request) {
	if h.GrpcConn == nil {
		http.Error(w, "no gprc connection", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	userId := chi.URLParam(r, "id")

	var reqDTO UpdatePatrolInfoRequestDTO
	reqDTO.UserId = userId

	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received put patrol information update request: %v", reqDTO)

	req := reqDTO.ToProto()
	log.Printf("Converted to put patrol information ptoro update request: %v", req)

	resp, err := h.GrpcClient.PutPatrolInfo(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !resp.Success {
		http.Error(w, ptr.DeferOrZero(resp.Message), http.StatusBadRequest)
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *PatrolBrokerHandler) PatchPatrolInfo(w http.ResponseWriter, r *http.Request) {
	if h.GrpcConn == nil {
		http.Error(w, "no gprc connection", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	userId := chi.URLParam(r, "id")

	var reqDTO UpdatePatrolInfoRequestDTO
	reqDTO.UserId = userId
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := reqDTO.ToProto()

	resp, err := h.GrpcClient.PatchPatrolInfo(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !resp.Success {
		http.Error(w, ptr.DeferOrZero(resp.Message), http.StatusBadRequest)
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *PatrolBrokerHandler) UpdatePatrolLocation(w http.ResponseWriter, r *http.Request) {
	if h.GrpcConn == nil {
		http.Error(w, "no gprc connection", http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	claims, ok := authmiddleware.GetClaims(r)
	if !ok {
		http.Error(w, "failed to obtain claims from auto token", http.StatusBadRequest)
	}

	var reqDTO UpdatePatrolLocationRequestDTO

	for _, role := range claims.Roles {
		switch role {
		case "ADMIN", "DISPATCHER":
			reqDTO.UserID = chi.URLParam(r, "id")
		case "PATROL":
			reqDTO.UserID = claims.Subject
		}
	}

	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req, err := reqDTO.ToProto()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.GrpcClient.UpdatePatrolLocation(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !resp.Success {
		http.Error(w, "TODO: implement message here - never will trigger for now", http.StatusBadRequest)
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *PatrolBrokerHandler) AssignPatrolToCrime(w http.ResponseWriter, r *http.Request) {
	if h.GrpcConn == nil {
		http.Error(w, "no gprc connection", http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var reqDTO AssignPatrolToCrimeRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received assign patrol to crime request: %v", reqDTO)

	req, err := reqDTO.ToProto()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Converted to assign patrol to crime proto request: %v", req)

	resp, err := h.GrpcClient.AssignPatrolToCrime(ctx, req)
	log.Printf("Assigned patrol to crime response: %v", resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !resp.Success {
		http.Error(w, ptr.DeferOrZero(resp.Message), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}
func (h *PatrolBrokerHandler) UpdatePatrolStatus(w http.ResponseWriter, r *http.Request) {
	if h.GrpcConn == nil {
		http.Error(w, "no gprc connection", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	claims, ok := authmiddleware.GetClaims(r)
	if !ok {
		http.Error(w, "failed to obtain claims from auto token", http.StatusBadRequest)
	}

	var reqDTO UpdatePatrolStatusRequestDTO

	for _, role := range claims.Roles {
		switch role {
		case "ADMIN", "DISPATCHER":
			reqDTO.PatrolId = chi.URLParam(r, "id")
		case "PATROL":
			reqDTO.PatrolId = claims.Subject
		}
	}

	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req, err := reqDTO.ToProto()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.GrpcClient.UpdatePatrolStatus(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !resp.Success {
		http.Error(w, ptr.DeferOrZero(resp.Message), http.StatusBadRequest)
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}
