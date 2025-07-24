package handler

import (
	"context"
	"crimeServiceApp/data"
	"crimeServiceApp/mapper"
	"crimeServiceApp/proto/crimepb"
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
func (s *CrimeServer) UpdateCrime(ctx context.Context, req *crimepb.UpdateCrimeReportRequest) (*crimepb.CrimeResponse, error) {
	panic("Not implemented")
}
func (s *CrimeServer) DeleteCrime(ctx context.Context, req *crimepb.DeleteCrimeRequest) (*crimepb.CrimeResponse, error) {
	panic("Not implemented")
}
