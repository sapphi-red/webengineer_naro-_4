package main

import (
	"net/http"
	"github.com/labstack/echo"
)

type ResponseData struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func createRoutes(e *echo.Echo) {
	withLogin := e.Group("")
	withLogin.Use(checkLogin)

	withLogin.GET("/cities/:cityName", getCityInfoHandler)
	withLogin.POST("/cities", postCityInfoHandler)
	withLogin.DELETE("/cities/:cityName", deleteCityInfoHandler)
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
