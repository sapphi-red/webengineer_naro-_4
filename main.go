package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Country struct {
	Code           string  `json:"code,omitempty"  db:"Code"`
	Name           string  `json:"name,omitempty"  db:"Name"`
	Continent      string  `json:"continent,omitempty"  db:"Continent"`
	Region         string  `json:"region,omitempty"  db:"Region"`
	SurfaceArea    float64 `json:"surfaceArea,omitempty"  db:"SurfaceArea"`
	IndepYear      int     `json:"indepYear,omitempty"  db:"IndepYear"`
	Population     int     `json:"population,omitempty"  db:"Population"`
	LifeExpectancy float64 `json:"lifeExpectancy,omitempty"  db:"LifeExpectancy"`
	GNP            float64 `json:"gnp,omitempty"  db:"GNP"`
	GNPOld         float64 `json:"gnpOld,omitempty"  db:"GNPOld"`
	LocalName      string  `json:"localName,omitempty"  db:"LocalName"`
	GovernmentForm string  `json:"governmentForm,omitempty"  db:"GovernmentForm"`
	HeadOfState    string  `json:"headOfState,omitempty"  db:"HeadOfState"`
	Capital       int     `json:"capital,omitempty"  db:"Capital"`
	Code2          string  `json:"code2,omitempty"  db:"Code2"`
}

type City struct {
	ID          int    `json:"id,omitempty" db:"ID"`
	Name        string `json:"name,omitempty" db:"Name"`
	CountryCode string `json:"countryCode,omitempty" db:"CountryCode"`
	District    string `json:"district,omitempty" db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

func main() {
	db, err := sqlx.Connect("mysql", fmt.Sprintf(
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
	fmt.Println("Connected!")

	cities := []City{}
	db.Select(&cities, "SELECT * FROM city WHERE CountryCode='JPN'")

	fmt.Println("日本の都市一覧")
	for _, city := range cities {
		fmt.Printf(
			"都市名: %s, 人口: %d人\n",
			city.Name,
			city.Population,	
		)
	}
	
}
