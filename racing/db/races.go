package db

import (
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error)
	// GetByID will return a single race based on the ID provided
	GetByID(id int64) (*racing.Race, error)
	// List Sports will return a list of sports
	ListSports(filter *racing.ListSportsRequestFilter) ([]*racing.SportEvent, error)
	// GetSportByID will return a single sport event based on the ID provided
	GetSportEventByID(id int64) (*racing.SportEvent, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

// Compiles the List of sports and applies filters if present
func (s *racesRepo) ListSports(filter *racing.ListSportsRequestFilter) ([]*racing.SportEvent, error) {
	var (
		err        error
		query      string
		args       []interface{}
		validField bool
	)
	// Create a bucket of valid fields that we can order by
	validFields := []string{"name", "id", "sport", "current_score", "advertised_start_time"}

	query = getSportQueries()[sportsList]
	query, args = s.applySportsFilter(query, filter)

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
	sports, err := s.scanSportEvents(rows)
	if err != nil {
		return nil, err
	}

	return sports, nil
}

// Get a sport by its Id
func (r *racesRepo) GetSportEventByID(id int64) (*racing.SportEvent, error) {
	// SQL Query to retrieve the race by its ID
	query := "SELECT id, name, advertised_start_time, sport, current_score  FROM races WHERE id = ?"

	// Execute query
	row := r.db.QueryRow(query, id)

	// Scan the row and get the sport event
	var sport racing.SportEvent
	var advertisedStart time.Time
	err := row.Scan(&sport.Id, &sport.Name, &advertisedStart, &sport.Sport, &sport.CurrentScore)
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

// Get a race by its Id
func (r *racesRepo) GetByID(id int64) (*racing.Race, error) {
	// SQL Query to retrieve the race by its ID
	query := "SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races WHERE id = ?"

	// Execute query
	row := r.db.QueryRow(query, id)

	// Scan the row and get race
	var race racing.Race
	var advertisedStart time.Time
	err := row.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // race not found
		}
		return nil, err
	}
	// Convert advertised start time to protobuf Timestamp
	ts, err := ptypes.TimestampProto(advertisedStart)
	if err != nil {
		return nil, err
	}
	race.AdvertisedStartTime = ts

	// Update status based on advertised start time
	if advertisedStart.Before(time.Now()) {
		race.Status = "CLOSED"
	} else {
		race.Status = "OPEN"
	}
	return &race, nil
}

// Compiles the List of races and applies filters if present
func (r *racesRepo) List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error) {
	var (
		err        error
		query      string
		args       []interface{}
		validField bool
	)
	// Create a bucket of valid fields that we can order by
	validFields := []string{"name", "number", "id", "meeting_id", "visible", "advertised_start_time"}

	query = getRaceQueries()[racesList]
	query, args = r.applyFilter(query, filter)

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

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	races, err := r.scanRaces(rows)
	if err != nil {
		return nil, err
	}
	// Update status based on advertised start time
	for _, race := range races {
		advertisedStart := time.Unix(race.AdvertisedStartTime.Seconds, int64(race.AdvertisedStartTime.Nanos))
		if advertisedStart.Before(time.Now()) {
			race.Status = "CLOSED"
		} else {
			race.Status = "OPEN"
		}
	}

	return races, nil
}

// Applies filters for sports and returns a SQL query
func (s *racesRepo) applySportsFilter(query string, filter *racing.ListSportsRequestFilter) (string, []interface{}) {
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

// Applies filters and returns a SQL query
func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}
	// Filters via Meeting ID - int array
	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}
	// Filters out races with visible set to false
	if filter.VisibleOnly {
		clauses = append(clauses, "visible = ?")
		args = append(args, true)
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

// Scans the SQL database and returns sport events
func (m *racesRepo) scanSportEvents(
	rows *sql.Rows,
) ([]*racing.SportEvent, error) {
	var sportEvents []*racing.SportEvent

	for rows.Next() {
		var sport racing.SportEvent
		var advertisedStart time.Time

		if err := rows.Scan(&sport.Id, &sport.Name, &advertisedStart, &sport.Sport, &sport.CurrentScore); err != nil {
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

// Scans the SQL database and returns races
func (m *racesRepo) scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		race.AdvertisedStartTime = ts

		races = append(races, &race)
	}

	return races, nil
}
