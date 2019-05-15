package router

import (
	"net/http"
	"github.com/labstack/echo"

	"github.com/sapphi-red/webengineer_naro-_4/database"
	"github.com/sapphi-red/webengineer_naro-_4/util"
)

func CreateRoutes(e *echo.Group) {
	e.GET("/cities/:cityName", getCityInfoHandler)
	e.POST("/cities", postCityInfoHandler)
	e.DELETE("/cities/:cityName", deleteCityInfoHandler)
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
	return util.ReturnSucessJSON(c)
}

func deleteCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")

	err := database.DeleteCity(cityName)
	if err != nil {
		return util.ReturnErrorJSON(c, "SQL Error")
	}
	return util.ReturnSucessJSON(c)
}
