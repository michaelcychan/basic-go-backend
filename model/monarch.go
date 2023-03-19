package model

// there must NOT be any space between json: and the "keyName"
type FullMonarchJson struct {
	Name        string `json:"name"`
	YearOfBirth int    `json:"birth_year"`
	YearOfDeath *int   `json:"death_year"`
	ReignStart  int    `json:"reign_start"`
	ReignEnd    *int   `json:"reign_end"`
}

type Room struct {
	Name string
	Temp int
}
