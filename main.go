package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Country struct {
	Code           string  `json:"code,omitempty"  db:"Code"`
	Name           string  `json:"name,omitempty"  db:"Name"`
	Continent      string  `json:"continent,omitempty"  db:"Continent"`
	Region         string  `json:"region,omitempty"  db:"Region"`
	SurfaceArea    float64 `json:"surface_area,omitempty"  db:"SurfaceArea"`
	IndepYear      int     `json:"indep_year,omitempty"  db:"IndepYear"`
	Population     int     `json:"population,omitempty"  db:"Population"`
	LifeExpectancy float64 `json:"life_expectancy,omitempty"  db:"LifeExpectancy"`
	GNP            float64 `json:"GNP,omitempty"  db:"GNP"`
	GNPOld         float64 `json:"GNP_old,omitempty"  db:"GNPOld"`
	LocalName      string  `json:"local_name,omitempty"  db:"LocalName"`
	GovernmentForm string  `json:"government_form,omitempty"  db:"GovernmentForm"`
	HeadOfState    string  `json:"head_of_state,omitempty"  db:"HeadOfState"`
	Capital        int     `json:"capital,omitempty"  db:"Capital"`
	Code2          string  `json:"code2,omitempty"  db:"Code2"`
}

type City struct {
	ID          int    `json:"ID,omitempty" db:"ID"`
	Name        string `json:"name,omitempty" db:"Name"`
	CountryCode string `json:"country_code,omitempty" db:"CountryCode"`
	District    string `json:"district,omitempty" db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

type ResponseData struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

var (
	db *sqlx.DB
)

func main() {
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

	store := createSessionStore()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.POST("/signup", postSignUpHandler)
	e.POST("/login", postLoginHandler)

	withLogin := e.Group("")
	withLogin.Use(checkLogin)

	withLogin.GET("/cities/:cityName", getCityInfoHandler)
	withLogin.POST("/cities", postCityInfoHandler)
	withLogin.DELETE("/cities/:cityName", deleteCityInfoHandler)

	e.Start(":12100")
}

func getCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")

	city := City{}
	db.Get(
		&city,
		`SELECT * FROM city WHERE Name = ?`,
		cityName,
	)
	if city.Name == "" {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, city)
}

func postCityInfoHandler(c echo.Context) error {
	data := new(City)
	err := c.Bind(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseData{
			Type:    "Error",
			Content: "Parse Error",
		})
	}

	_, err = db.Exec(
		`INSERT INTO city (Name, CountryCode, District, Population) VALUES (?, ?, ?, ?)`,
		data.Name,
		data.CountryCode,
		data.District,
		data.Population,
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseData{
			Type:    "Error",
			Content: "SQL Error",
		})
	}
	return c.JSON(http.StatusOK, ResponseData{
		Type:    "Success",
		Content: "Success",
	})
}

func deleteCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")

	_, err := db.Exec(
		`DELETE FROM city WHERE Name = ?`,
		cityName,
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseData{
			Type:    "Error",
			Content: "SQL Error",
		})
	}
	return c.JSON(http.StatusOK, ResponseData{
		Type:    "Success",
		Content: "Success",
	})
}