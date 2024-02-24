package db

import (
	"database/sql"
	"strings"
	"sync"
	"time"

	"command-line-arguments/home/james/Documents/GitHub/EntainGroupTest/sports/proto/sports/sports.pb.go"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"git.neds.sh/matty/entain/sports/proto/sports"
)

// SportsRepo provides repository access to sports.
type SportsRepo interface {
	// Init will initialise our sports repository.
	Init() error

	// List will return a list of sports.
	List(filter *sports.ListSportsRequestFilter) ([]*sports.SportEvent, error)
	// GetByID will return a single sport based on the ID provided
	GetByID(id int64) (*sports.SportEvent, error)
}

type sportsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewSportsRepo creates a new sports repository.
func NewSportsRepo(db *sql.DB) SportsRepo {
	return &sportsRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *sportsRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

// Get a race by its Id
func (r *sportsRepo) GetByID(id int64) (*sports.sportEvent, error) {
	// SQL Query to retrieve the race by its ID
	query := "SELECT id, name, advertised_start_time, sport, current_score  FROM races WHERE id = ?"

	// Execute query
	row := r.db.QueryRow(query, id)

	// Scan the row and get the sport event
	var sport sports.sportEvent
	var advertisedStart time.Time
	err := row.Scan(&sport.Id, &sport.Name, &advertisedStart, &sport.sport, &sport.current_score)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // race not found
		}
		return nil, err
	}
	// Convert advertised start time to protobuf Timestamp
	/*ts, err := ptypes.TimestampProto(advertisedStart)
	if err != nil {
		return nil, err
	}
	race.AdvertisedStartTime = ts

	// Update status based on advertised start time
	if advertisedStart.Before(time.Now()) {
		race.Status = "CLOSED"
	} else {
		race.Status = "OPEN"
	}*/
	return &sport, nil
}

// Compiles the List of sports and applies filters if present
func (s *sportsRepo) List(filter *sports.ListSportsRequestFilter) ([]*sports.sportEvent, error) {
	var (
		err        error
		query      string
		args       []interface{}
		validField bool
	)
	// Create a bucket of valid fields that we can order by
	validFields := []string{"name", "id", "sport", "current_score", "advertised_start_time"}

	query = getSportQueries()[sportsList]
	query, args = s.applyFilter(query, filter)

	// Check if orderBy is provided in the filter
	if filter != nil && filter.OrderBy != "" {
		// Check if orderBy is a valid field
		validField = false
		for _, field := range validFields {
			if filter.OrderBy == field {
				validField = true
				break
			}
		}
		// If orderBy is valid we order by that field else order by advertised_start_time by default
		if validField {
			query += " ORDER BY " + filter.OrderBy
		} else {
			query += " ORDER BY advertised_start_time"
		}
		// Check if sorting direction is provided
		if filter.Sort != "" {
			if strings.ToLower(filter.Sort) == "desc" {
				query += " DESC"
			} else {
				query += " ASC"
			}
		}
	} else {
		// If orderBy is not provided, order by advertised_start_time
		query += " ORDER BY advertised_start_time"
		// Default sorting direction
		query += " ASC"
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	sports, err := s.scanSports(rows)
	if err != nil {
		return nil, err
	}

	return sports, nil
}

// applies filters and returns a SQL query
func (s *sportsRepo) applyFilter(query string, filter *sports.ListSportsRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)
	var availableSports = []string{"basketball", "soccer", "hockey", "rugby league", "afl"}

	if filter == nil {
		return query, args
	}
	// Filters via ID's - int array
	if len(filter.Ids) > 0 {
		clauses = append(clauses, "id IN ("+strings.Repeat("?,", len(filter.Ids)-1)+"?)")

		for _, ID := range filter.Ids {
			args = append(args, ID)
		}
	}

	// Filter via sport
	if filter.Sport != "" {
		var validSport bool
		for _, sport := range availableSports {
			if strings.EqualFold(filter.Sport, sport) {
				validSport = true
				break
			}
			if validSport {
				clauses = append(clauses, "sport = ?")
				args = append(args, filter.Sport)
			}
		}
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

// scans the SQL database and returns races
func (m *sportsRepo) scanSports(
	rows *sql.Rows,
) ([]*sports.sportEvent, error) {
	var sportEvents []*sports.sportEvent

	for rows.Next() {
		var sport sports.sportEvent
		var advertisedStart time.Time

		if err := rows.Scan(&sport.Id, &sport.Name, &advertisedStart, &sport.sport, &sport.current_score); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		sport.AdvertisedStartTime = ts

		sportEvents = append(sportEvents, &sport)
	}

	return sportEvents, nil
}
