package service

import (
	"github.com/sibeyzoran/EntainGroupTest/racing/db"
	"github.com/sibeyzoran/EntainGroupTest/racing/proto/sports"

	"golang.org/x/net/context"
)

type Sporting interface {
	// ListSports will return a collection of sports.
	ListSports(ctx context.Context, in *sports.ListSportsRequest) (*sports.ListSportsResponse, error)
	// GetRaceByID will return a single sport
	GetSportByID(ctx context.Context, in *sports.GetSportByIDRequest) (*sports.GetSportByIDResponse, error)
}

// sportingService implements the Sporting interface.
type sportingService struct {
	racesRepo db.RacesRepo
}

// NewSportingService instantiates and returns a new sportingService.
func NewSportingService(racesRepo db.RacesRepo) Sporting {
	return &sportingService{racesRepo}
}

func (s *sportingService) ListSports(ctx context.Context, in *sports.ListSportsRequest) (*sports.ListSportsResponse, error) {
	sportEvents, err := s.racesRepo.ListSports(in.Filter)
	if err != nil {
		return nil, err
	}
	// Create a new ListSportsResponse (unsure why I had to make this into a variable)
	response := &sports.ListSportsResponse{Sports: sportEvents}
	return response, nil
}

func (s *sportingService) GetSportByID(ctx context.Context, in *sports.GetSportByIDRequest) (*sports.GetSportByIDResponse, error) {
	sport, err := s.racesRepo.GetSportEventByID(in.Id)
	if err != nil {
		return nil, err
	}

	return &sports.GetSportByIDResponse{Sport: sport}, nil
}
