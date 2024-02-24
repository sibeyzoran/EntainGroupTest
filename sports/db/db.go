package db

import (
	"fmt"
	"math/rand"
	"time"

	"syreclabs.com/go/faker"
)

// Possible Sports for seeding into DB Table
var sports = []string{"Basketball", "Soccer", "Hockey", "Rugby League", "AFL"}

func (s *sportsRepo) seed() error {
	// Prepare SQL table if it doesn't exist
	statement, err := s.db.Prepare(`CREATE TABLE IF NOT EXISTS sports (id INTEGER PRIMARY KEY, name TEXT, advertised_start_time DATETIME, sport TEXT, current_score TEXT)`)
	if err == nil {
		_, err = statement.Exec()
	}

	// Insert fake data into the table
	for i := 1; i <= 100; i++ {
		// Make a random team match up
		teamA := faker.Team().Name()
		teamB := faker.Team().Name()
		name := fmt.Sprintf("%s VS %s", teamA, teamB)

		// Select a random sport
		sportIndex := rand.Intn(len(sports))
		sport := sports[sportIndex]

		// Make a random score
		currentScore := fmt.Sprintf("%d-%d", rand.Intn(151), rand.Intn(151))

		statement, err = s.db.Prepare(`INSERT OR IGNORE INTO races(id, name, advertised_start_time, sport, current_score) VALUES (?,?,?,?,?)`)
		if err == nil {
			_, err = statement.Exec(
				i,
				name,
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)).Format(time.RFC3339),
				sport,
				currentScore,
			)
		}
	}

	return err
}
