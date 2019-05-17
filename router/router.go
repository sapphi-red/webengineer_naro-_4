package router

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

func CreateRoutes(e *echo.Echo, db *sqlx.DB) {
	g := e.Group("")
	g.Use(checkLogin)

	g.GET("/whoami", getWhoAmIHandler)

	g.GET("/cities/:cityName", getCityInfoHandler)
	g.POST("/cities", postCityInfoHandler)
	g.DELETE("/cities/:cityName", deleteCityInfoHandler)

	g.GET("/countries", getCountriesInfoHandler)
	g.GET("/countries/:countryName", getCountryInfoHandler)
}
