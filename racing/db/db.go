package db

import (
	"fmt"
	"math/rand"
	"time"

	"syreclabs.com/go/faker"
)

// Possible Sports for seeding into DB Table
var sportsData = []string{"Basketball", "Soccer", "Hockey", "Rugby League", "AFL"}

func (r *racesRepo) seed() error {
	// Prepare SQL table if it doesn't exist
	statement, err := r.db.Prepare(`CREATE TABLE IF NOT EXISTS races (id INTEGER PRIMARY KEY, meeting_id INTEGER, name TEXT, number INTEGER, visible INTEGER, advertised_start_time DATETIME)`)
	if err == nil {
		_, err = statement.Exec()
	}

	// Populate with fake data
	for i := 1; i <= 100; i++ {
		statement, err = r.db.Prepare(`INSERT OR IGNORE INTO races(id, meeting_id, name, number, visible, advertised_start_time) VALUES (?,?,?,?,?,?)`)
		if err == nil {
			_, err = statement.Exec(
				i,
				faker.Number().Between(1, 10),
				faker.Team().Name(),
				faker.Number().Between(1, 12),
				faker.Number().Between(0, 1),
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)).Format(time.RFC3339),
			)
		}
	}

	// Prepare sports SQL table if it doesn't exist
	sportStatement, err := r.db.Prepare(`CREATE TABLE IF NOT EXISTS sports (id INTEGER PRIMARY KEY, name TEXT, advertised_start_time DATETIME, sport TEXT, current_score TEXT)`)
	if err == nil {
		_, err = sportStatement.Exec()
	}

	// Insert fake data into the table
	for i := 1; i <= 100; i++ {
		// Make a random team match up
		teamA := faker.Team().Name()
		teamB := faker.Team().Name()
		name := fmt.Sprintf("%s VS %s", teamA, teamB)

		// Select a random sport
		sportIndex := rand.Intn(len(sportsData))
		sport := sportsData[sportIndex]

		// Make a random score
		currentScore := fmt.Sprintf("%d-%d", rand.Intn(151), rand.Intn(151))

		statement, err = r.db.Prepare(`INSERT OR IGNORE INTO sports(id, name, advertised_start_time, sport, current_score) VALUES (?,?,?,?,?)`)
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
