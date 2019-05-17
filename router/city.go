package router

import (
	"github.com/labstack/echo"
	"github.com/sapphi-red/webengineer_naro-_4/database"
	"net/http"
)

func getCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")

	city, _ := database.GetCity(cityName)
	if city.Name == "" {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, city)
}

func postCityInfoHandler(c echo.Context) error {
	data := new(database.City)
	err := c.Bind(data)
	if err != nil {
		return returnErrorJSON(c, "Parse Error")
	}

	err = database.AddCity(data)
	if err != nil {
		return returnErrorJSON(c, "SQL Error")
	}
	return returnSuccessJSON(c)
}

func deleteCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")

	err := database.DeleteCity(cityName)
	if err != nil {
		return returnErrorJSON(c, "SQL Error")
	}
	return returnSuccessJSON(c)
}
