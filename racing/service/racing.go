package service

import (
	"github.com/sibeyzoran/EntainGroupTest/racing/db"
	"github.com/sibeyzoran/EntainGroupTest/racing/proto/racing"

	"golang.org/x/net/context"
)

type Racing interface {
	// ListRaces will return a collection of races.
	ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error)
	// GetRaceByID will return a single race
	GetRaceByID(ctx context.Context, in *racing.GetRaceByIDRequest) (*racing.GetRaceByIDResponse, error)
}

// racingService implements the Racing interface.
type racingService struct {
	racesRepo db.RacesRepo
}

// NewRacingService instantiates and returns a new racingService.
func NewRacingService(racesRepo db.RacesRepo) Racing {
	return &racingService{racesRepo}
}

// List all races
func (r *racingService) ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error) {
	races, err := r.racesRepo.List(in.Filter)
	if err != nil {
		return nil, err
	}

	return &racing.ListRacesResponse{Races: races}, nil
}

// Gets and returns a single race
func (r *racingService) GetRaceByID(ctx context.Context, in *racing.GetRaceByIDRequest) (*racing.GetRaceByIDResponse, error) {
	race, err := r.racesRepo.GetByID(in.Id)
	if err != nil {
		return nil, err
	}
	return &racing.GetRaceByIDResponse{Race: race}, nil
}
