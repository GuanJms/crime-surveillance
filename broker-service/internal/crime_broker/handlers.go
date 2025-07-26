package crimebroker

import (
	"brokerServiceApp/internal/authmiddleware"
	"brokerServiceApp/internal/crime_broker/proto/crimepb"
	"brokerServiceApp/utils"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var secret []byte

func initializeSecret() {
	secretStr := os.Getenv("SECRET")
	if secretStr == "" {
		log.Fatal("SECRET environment variable is required")
	}
	secret = []byte(secretStr)
}

type CrimeBrokerHandler struct {
	GrpcClient crimepb.CrimeServiceClient
	GrpcConn   *grpc.ClientConn
}

// TODO: Add graceful shutdown managemetn of conn
func NewCrimeBrokerHandler() (*CrimeBrokerHandler, error) {
	initializeSecret()
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

// TODO: implement list all crimes handlers
func (cb *CrimeBrokerHandler) ListAllCrimes(w http.ResponseWriter, r *http.Request) {
	if cb.GrpcConn == nil {
		utils.ErrorJSON(w, errors.New("crime broker handler has no grpc connection"), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	getCrimeRequest := crimepb.GetCrimesRequest{}

	resp, err := cb.GrpcClient.GetAllCrimes(ctx, &getCrimeRequest)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	crimesDTOs := CrimesProtoToDTOs(resp.Crimes)
	utils.WriteJSON(w, http.StatusOK, crimesDTOs)
}

// TODO: implement submit new crime
func (cb *CrimeBrokerHandler) SubmitNewCrime(w http.ResponseWriter, r *http.Request) {
	if cb.GrpcConn == nil {
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

	resp, err := cb.GrpcClient.SubmitNewCrimeReport(ctx, &req)
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
func (cb *CrimeBrokerHandler) PutCrime(w http.ResponseWriter, r *http.Request) {
	crimeID := chi.URLParam(r, "id")
	//reporterId should not use authorization token since - only admin/officer will update it

	var reqDTO UpdateCrimeReportRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	req, err := reqDTO.toProto()
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	req.Id = crimeID

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := cb.GrpcClient.PutCrime(ctx, req)
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

func (cb *CrimeBrokerHandler) PatchCrime(w http.ResponseWriter, r *http.Request) {
	crimeId := chi.URLParam(r, "id")
	//reporterId should not use authorization token since - only admin/officer will update it

	var reqDTO UpdateCrimeReportRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		http.Error(w, "Invalid JSON during decoding request", http.StatusBadRequest)
		return
	}
	log.Printf("Received update crime report request DTO - %v", reqDTO)

	req, err := reqDTO.toProto()
	if err != nil {
		http.Error(w, "error creating reqDTO to Proto class", http.StatusBadRequest)
		return
	}
	req.Id = crimeId

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := cb.GrpcClient.PatchCrime(ctx, req)
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

func (cb *CrimeBrokerHandler) DeleteCrime(w http.ResponseWriter, r *http.Request) {
	var deleteRequest crimepb.DeleteCrimeRequest

	deleteRequest.Id = chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := cb.GrpcClient.DeleteCrime(ctx, &deleteRequest)
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

func (cb *CrimeBrokerHandler) AddTo(r chi.Router) {
	r.Route("/crimes", func(crime chi.Router) {
		crime.Get("/", cb.ListAllCrimes)
		crime.Group(func(protected chi.Router) {
			protected.Use(authmiddleware.JWTMiddleware(secret))
			protected.Post("/", cb.SubmitNewCrime)
			protected.With(authmiddleware.RequireRole("PATROL", "DISPATCHER", "ADMIN")).Put("/{id}", cb.PutCrime)
			protected.With(authmiddleware.RequireRole("PATROL", "DISPATCHER", "ADMIN")).Patch("/{id}", cb.PatchCrime)
			protected.With(authmiddleware.RequireRole("ADMIN")).Delete("/{id}", cb.DeleteCrime)
		})
	})
}
