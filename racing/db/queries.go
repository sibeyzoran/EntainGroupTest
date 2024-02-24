package db

const (
	racesList  = "list"
	sportsList = "list"
)

func getRaceQueries() map[string]string {
	return map[string]string{
		racesList: `
			SELECT 
				id, 
				meeting_id, 
				name, 
				number, 
				visible, 
				advertised_start_time 
			FROM races
		`,
	}
}

func getSportQueries() map[string]string {
	return map[string]string{
		sportsList: `
			SELECT 
				id, 
				name, 
				advertised_start_time , 
				sport, 
				current_score
			FROM sports
		`,
	}
}
