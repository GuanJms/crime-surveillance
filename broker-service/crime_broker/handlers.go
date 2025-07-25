package crimebroker

import (
	"brokerServiceApp/crime_broker/proto/crimepb"
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
	// TODO: implement auth for tracking the user_id
	// req.ReporterId = r.Context().Value("user_id").(string)
	// req.ReporterId = uuid.New().String()

	// TODO: delete dev reporter_id
	req.ReporterId = "c2473b3c-ab82-4ed4-a814-1f5dc29e730f"

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
	crimeID := chi.URLParam(r, "id")
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
	req.Id = crimeID

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

// TODO: implement delete crime
func (cb *CrimeBrokerHandler) DeleteCrime(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (cb *CrimeBrokerHandler) AddTo(mux *chi.Mux) {
	mux.Get("/crimes", cb.ListAllCrimes)
	mux.Post("/crimes", cb.SubmitNewCrime)
	mux.Put("/crimes/{id}", cb.PutCrime)
	mux.Patch("/crimes/{id}", cb.PatchCrime)
	mux.Delete("/crimes/{id}", cb.DeleteCrime)
}
