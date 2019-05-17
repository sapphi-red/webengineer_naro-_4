package database

type Country struct {
	Code           string   `json:"code,omitempty"  db:"Code"`
	Name           string   `json:"name,omitempty"  db:"Name"`
	Continent      string   `json:"continent,omitempty"  db:"Continent"`
	Region         string   `json:"region,omitempty"  db:"Region"`
	SurfaceArea    float64  `json:"surface_area,omitempty"  db:"SurfaceArea"`
	IndepYear      *int     `json:"indep_year,omitempty"  db:"IndepYear"`
	Population     int      `json:"population,omitempty"  db:"Population"`
	LifeExpectancy *float64 `json:"life_expectancy,omitempty"  db:"LifeExpectancy"`
	GNP            *float64 `json:"GNP,omitempty"  db:"GNP"`
	GNPOld         *float64 `json:"GNP_old,omitempty"  db:"GNPOld"`
	LocalName      string   `json:"local_name,omitempty"  db:"LocalName"`
	GovernmentForm string   `json:"government_form,omitempty"  db:"GovernmentForm"`
	HeadOfState    *string  `json:"head_of_state,omitempty"  db:"HeadOfState"`
	Capital        *int     `json:"capital,omitempty"  db:"Capital"`
	Code2          string   `json:"code2,omitempty"  db:"Code2"`
}

type CountryCities struct {
	Country string `json:"country,omitempty"  db:"Country"`
	Name    string `json:"name,omitempty"  db:"Name"`
}

func GetCountries() ([]Country, error) {
	countries := []Country{}
	err := db.Select(
		&countries,
		`SELECT Name FROM country`,
	)
	return countries, err
}

func GetCountryCities(name string) ([]CountryCities, error) {
	countryCities := []CountryCities{}
	err := db.Select(
		&countryCities,
		`SELECT city.Name FROM country JOIN city ON country.Code=city.CountryCode WHERE country.Name = ?`,
		name,
	)
	return countryCities, err
}
