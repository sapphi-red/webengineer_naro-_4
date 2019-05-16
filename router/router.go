package router

import (
	"github.com/labstack/echo"
	"net/http"

	"github.com/sapphi-red/webengineer_naro-_4/database"
	"github.com/sapphi-red/webengineer_naro-_4/util"
)

type Me struct {
	Username string `json:"username,omitempty"  db:"username"`
}

func CreateRoutes(e *echo.Group) {
	e.GET("/whoami", getWhoAmIHandler)

	e.GET("/cities/:cityName", getCityInfoHandler)
	e.POST("/cities", postCityInfoHandler)
	e.DELETE("/cities/:cityName", deleteCityInfoHandler)

	e.GET("/countries", getCountriesInfoHandler)
	e.GET("/countries/:countryName", getCountryInfoHandler)
}

func getWhoAmIHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, Me{
		Username: c.Get("userName").(string),
	})
}

func getCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")

	city := database.GetCity(cityName)
	if city.Name == "" {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, city)
}

func postCityInfoHandler(c echo.Context) error {
	data := new(database.City)
	err := c.Bind(data)
	if err != nil {
		return util.ReturnErrorJSON(c, "Parse Error")
	}

	err = database.AddCity(data)
	if err != nil {
		return util.ReturnErrorJSON(c, "SQL Error")
	}
	return util.ReturnSuccessJSON(c)
}

func deleteCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")

	err := database.DeleteCity(cityName)
	if err != nil {
		return util.ReturnErrorJSON(c, "SQL Error")
	}
	return util.ReturnSuccessJSON(c)
}

func getCountriesInfoHandler(c echo.Context) error {
	countries := database.GetCountries()
	return c.JSON(http.StatusOK, countries)
}

func getCountryInfoHandler(c echo.Context) error {
	countryName := c.Param("countryName")
	country := database.GetCountryCities(countryName)
	return c.JSON(http.StatusOK, country)
}
