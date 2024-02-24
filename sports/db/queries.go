package db

const (
	sportsList = "list"
)

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
