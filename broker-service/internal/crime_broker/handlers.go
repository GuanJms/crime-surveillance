package crimebroker

import (
	"brokerServiceApp/internal/authmiddleware"
	"brokerServiceApp/internal/crime_broker/proto/crimepb"
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

type CrimeBrokerHandler struct {
	GrpcClient crimepb.CrimeServiceClient
	GrpcConn   *grpc.ClientConn
}

// TODO: Add graceful shutdown managemetn of conn
func NewCrimeBrokerHandler() (*CrimeBrokerHandler, error) {
	conn, err := grpc.NewClient("crime-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := crimepb.NewCrimeServiceClient(conn)
	return &CrimeBrokerHandler{
		GrpcClient: client,
		GrpcConn:   conn,
	}, nil
}

func (h *CrimeBrokerHandler) AddTo(r chi.Router) {
	r.Route("/crimes", func(crime chi.Router) {
		crime.Get("/", h.ListAllCrimes)
		crime.Group(func(protected chi.Router) {
			protected.Use(authmiddleware.JWTMiddleware(utils.Secret))
			protected.Post("/", h.SubmitNewCrime)
			protected.With(authmiddleware.RequireRole("PATROL", "DISPATCHER", "ADMIN")).Put("/{id}", h.PutCrime)
			protected.With(authmiddleware.RequireRole("PATROL", "DISPATCHER", "ADMIN")).Patch("/{id}", h.PatchCrime)
			protected.With(authmiddleware.RequireRole("PATROL", "DISPATCHER", "ADMIN")).Delete("/{id}", h.DeleteCrime)
		})
	})
}

// TODO: implement list all crimes handlers
func (h *CrimeBrokerHandler) ListAllCrimes(w http.ResponseWriter, r *http.Request) {
	if h.GrpcConn == nil {
		utils.ErrorJSON(w, errors.New("crime broker handler has no grpc connection"), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	getCrimeRequest := crimepb.GetCrimesRequest{}

	resp, err := h.GrpcClient.GetAllCrimes(ctx, &getCrimeRequest)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	crimesDTOs := CrimesProtoToDTOs(resp.Crimes)
	utils.WriteJSON(w, http.StatusOK, crimesDTOs)
}

// TODO: implement submit new crime
func (h *CrimeBrokerHandler) SubmitNewCrime(w http.ResponseWriter, r *http.Request) {
	if h.GrpcConn == nil {
		utils.ErrorJSON(w, errors.New("crime broker handler has no grpc connection"), http.StatusInternalServerError)
		return
	}
	var req crimepb.CrimeReportRequest
	claims, ok := authmiddleware.GetClaims(r)
	if !ok {
		utils.ErrorJSON(w, errors.New("failed to achieve user_id from the claims"), http.StatusInternalServerError)
		return
	}

	// access reporter_id through subject in registeredClaims
	req.ReporterId, _ = claims.GetSubject()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := h.GrpcClient.SubmitNewCrimeReport(ctx, &req)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// crime response in resp including id, successful, message
	if !resp.Successful {
		utils.ErrorJSON(w, errors.New(resp.Message), http.StatusBadRequest)
		return
	}

	// successful response
	utils.WriteJSON(w, http.StatusOK, resp)
}

// TODO: implement update crime
func (h *CrimeBrokerHandler) PutCrime(w http.ResponseWriter, r *http.Request) {
	crimeID := chi.URLParam(r, "id")
	//reporterId should not use authorization token since - only admin/officer will update it

	var reqDTO UpdateCrimeReportRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	req, err := reqDTO.toProto()
	if err != nil {
		http.Error(w, "Invalid JSON proto", http.StatusBadRequest)
		return
	}
	req.Id = crimeID

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := h.GrpcClient.PutCrime(ctx, req)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	// crime response in resp including id, successful, message
	if !resp.Successful {
		utils.ErrorJSON(w, errors.New(resp.Message), http.StatusBadRequest)
		return
	}
	// successful response
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *CrimeBrokerHandler) PatchCrime(w http.ResponseWriter, r *http.Request) {
	crimeId := chi.URLParam(r, "id")
	//reporterId should not use authorization token since - only admin/officer will update it

	var reqDTO UpdateCrimeReportRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		http.Error(w, "Invalid JSON during decoding request", http.StatusBadRequest)
		return
	}
	// log.Printf("Received update crime report request DTO - %v", reqDTO)

	req, err := reqDTO.toProto()
	if err != nil {
		http.Error(w, "error creating reqDTO to Proto class", http.StatusBadRequest)
		return
	}
	req.Id = crimeId

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := h.GrpcClient.PatchCrime(ctx, req)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	// crime response in resp including id, successful, message
	if !resp.Successful {
		utils.ErrorJSON(w, errors.New(resp.Message), http.StatusBadRequest)
		return
	}
	// successful response
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *CrimeBrokerHandler) DeleteCrime(w http.ResponseWriter, r *http.Request) {
	var deleteRequest crimepb.DeleteCrimeRequest

	deleteRequest.Id = chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := h.GrpcClient.DeleteCrime(ctx, &deleteRequest)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	// crime response in resp including id, successful, message
	if !resp.Successful {
		utils.ErrorJSON(w, errors.New(resp.Message), http.StatusBadRequest)
		return
	}
	// successful response
	utils.WriteJSON(w, http.StatusOK, resp)
}
