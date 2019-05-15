package main

import (
	"fmt"
	"net/http"
	"github.com/labstack/echo"
)

func return400(c echo.Context, err error) error {
	return c.String(http.StatusBadRequest, err.Error())
}

func return500(c echo.Context, err error) error {
	return c.String(http.StatusInternalServerError, err.Error())
}

func returnDBError(c echo.Context, err error) error {
	return return500(c, fmt.Errorf("db error: %v", err))
}

func returnErrorJSON(c echo.Context, content string) error {
	return c.JSON(http.StatusBadRequest, ResponseData{
		Type:    "Error",
		Content: content,
	})
}

func returnSucessJSON(c echo.Context) error {
	return c.JSON(http.StatusOK, ResponseData{
		Type:    "Success",
		Content: "Success",
	})
}