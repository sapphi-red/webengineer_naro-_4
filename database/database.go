package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

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

type City struct {
	ID          int    `json:"ID,omitempty" db:"ID"`
	Name        string `json:"name,omitempty" db:"Name"`
	CountryCode string `json:"country_code,omitempty" db:"CountryCode"`
	District    string `json:"district,omitempty" db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

type CountryCities struct {
	Country string `json:"country,omitempty"  db:"Country"`
	Name    string `json:"name,omitempty"  db:"Name"`
}

var (
	db *sqlx.DB
)

func ConnectDB() *sqlx.DB {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOSTNAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}

	db = _db

	return _db
}

func GetCity(name string) City {
	city := City{}
	db.Get(
		&city,
		`SELECT * FROM city WHERE Name = ?`,
		name,
	)
	return city
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

func GetCountries() []Country {
	countries := []Country{}
	db.Select(
		&countries,
		`SELECT * FROM country`,
	)
	return countries
}

func GetCountryCities(name string) []CountryCities {
	countryCities := []CountryCities{}
	db.Select(
		&countryCities,
		`SELECT country.Name AS Country, city.Name FROM country JOIN city ON country.Code=city.CountryCode WHERE country.Name = ?`,
		name,
	)
	return countryCities
}
