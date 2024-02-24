package service

import (
	"context"

	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
)

type Sports interface {
	// ListSports will return a collection of sports.
	ListSports(ctx context.Context, in *racing.ListSportsRequest) (*racing.ListSportsResponse, error)
	// GetSportByID will return a single sport.
	GetSportByID(ctx context.Context, in *racing.GetSportByIDRequest) (*racing.GetSportByIDResponse, error)
}

// sportsService implements the Sports interface.
type sportsService struct {
	sportsRepo db.SportsRepo
}

// NewSportsService instantiates and returns a new sportsService.
func NewSportsService(sportsRepo db.SportsRepo) Sports {
	return &sportsService{sportsRepo}
}

// Lists all sports
func (s *sportsService) ListSports(ctx context.Context, in *racing.ListSportsRequest) (*racing.ListSportsResponse, error) {
	sports, err := s.sportsRepo.List(in.Filter)
	if err != nil {
		return nil, err
	}
	return &racing.ListSportsResponse{Sports: sports}, nil
}

// Gets a single sport
func (s *sportsService) GetSportByID(ctx context.Context, in *racing.GetSportByIDRequest) (*racing.GetSportByIDResponse, error) {
	sport, err := s.sportsRepo.GetByID(in.Id)
	if err != nil {
		return nil, err
	}
	return &racing.GetSportByIDResponse{Sport: sport}, nil
}
