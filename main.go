package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type ResponseData struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

var (
	db *sqlx.DB
)

func main() {
	db = connectDB()
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

	city := getCity(cityName)
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

	err = addCity(data)
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

	err := deleteCity(cityName)
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