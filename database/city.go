package database

type City struct {
	ID          int    `json:"ID,omitempty" db:"ID"`
	Name        string `json:"name,omitempty" db:"Name"`
	CountryCode string `json:"country_code,omitempty" db:"CountryCode"`
	District    string `json:"district,omitempty" db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

func GetCity(name string) (*City, error) {
	city := City{}
	err := db.Get(
		&city,
		`SELECT * FROM city WHERE Name = ?`,
		name,
	)
	return &city, err
}

func AddCity(city *City) error {
	_, err := db.Exec(
		`INSERT INTO city (Name, CountryCode, District, Population) VALUES (?, ?, ?, ?)`,
		city.Name,
		city.CountryCode,
		city.District,
		city.Population,
	)
	return err
}

func DeleteCity(name string) error {
	_, err := db.Exec(
		`DELETE FROM city WHERE Name = ?`,
		name,
	)
	return err
}
