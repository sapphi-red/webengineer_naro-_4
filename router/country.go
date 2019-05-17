package router

import (
	"github.com/labstack/echo"
	"github.com/sapphi-red/webengineer_naro-_4/database"
	"net/http"
)

func getCountriesInfoHandler(c echo.Context) error {
	countries, _ := database.GetCountries()
	return c.JSON(http.StatusOK, countries)
}

func getCountryInfoHandler(c echo.Context) error {
	countryName := c.Param("countryName")
	country, _ := database.GetCountryCities(countryName)
	return c.JSON(http.StatusOK, country)
}
