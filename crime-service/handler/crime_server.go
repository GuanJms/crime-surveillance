package handler

import (
	"context"
	"crimeServiceApp/data"
	"crimeServiceApp/mapper"
	"crimeServiceApp/proto/crimepb"
	"log"
)

type CrimeServer struct {
	crimepb.UnimplementedCrimeServiceServer
	CrimeModels *data.CrimeModel
}

func (s *CrimeServer) GetAllCrimes(ctx context.Context, req *crimepb.GetCrimesRequest) (*crimepb.GetCrimesResponse, error) {
	crimes, err := s.CrimeModels.Repo.GetAllCrimes() // get all the crimes from db
	if err != nil {
		return nil, err
	}

	crimesProto := make([]*crimepb.Crime, 0)

	for _, c := range crimes {
		// take each crime and convert to crimeProto
		cp, err := mapper.ToProtoCrime(c)
		if err != nil {
			return nil, err
		}
		crimesProto = append(crimesProto, cp)
	}

	resp := &crimepb.GetCrimesResponse{
		Crimes: crimesProto,
	}
	return resp, nil
}
func (s *CrimeServer) SubmitNewCrimeReport(ctx context.Context, req *crimepb.CrimeReportRequest) (*crimepb.CrimeResponse, error) {
	newCrime, err := data.NewCrimeFromProto(req)
	if err != nil {
		return nil, err
	}
	err = s.CrimeModels.Repo.InsertNewCrime(newCrime)
	if err != nil {
		return nil, err
	}
	resp := &crimepb.CrimeResponse{
		Id:         newCrime.ID,
		Successful: true,
		Message:    "New crime successfully created",
	}
	return resp, nil
}
func (s *CrimeServer) PutCrime(ctx context.Context, req *crimepb.UpdateCrimeReportRequest) (*crimepb.CrimeResponse, error) {
	crime, err := data.NewCrimeFromProtoUpdateRequest(req)
	if err != nil {
		return nil, err
	}
	err = s.CrimeModels.Repo.PutCrime(crime)
	if err != nil {
		return nil, err
	}
	resp := &crimepb.CrimeResponse{
		Id:         crime.ID,
		Successful: true,
		Message:    "Successfully put executed",
	}
	return resp, nil

}
func (s *CrimeServer) PatchCrime(ctx context.Context, req *crimepb.UpdateCrimeReportRequest) (*crimepb.CrimeResponse, error) {
	update := data.NewCrimeUpdateFromProtoUpdateRequest(req)

	log.Printf("updating patch crime with update %v", update)

	err := s.CrimeModels.Repo.PatchCrime(update)
	if err != nil {
		return nil, err
	}
	resp := &crimepb.CrimeResponse{
		Id:         update.ID,
		Successful: true,
		Message:    "Successfully patch update executed",
	}
	return resp, nil
}
func (s *CrimeServer) DeleteCrime(ctx context.Context, req *crimepb.DeleteCrimeRequest) (*crimepb.CrimeResponse, error) {
	crimeID := req.Id
	err := s.CrimeModels.Repo.DeleteCrime(crimeID)
	if err != nil {
		return nil, err
	}

	resp := &crimepb.CrimeResponse{
		Id:         crimeID,
		Successful: true,
		Message:    "Successfully deleted crime record",
	}
	return resp, nil
}
