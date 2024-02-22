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

	return r.scanRaces(rows)
}

// applies filters and returns a SQL query
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

// scans the SQL database and returns races
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
