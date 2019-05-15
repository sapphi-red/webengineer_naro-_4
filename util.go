package main

import (
	"net/http"
	"github.com/labstack/echo"
)

func return400(c echo.Context, err error) error {
	return c.String(http.StatusBadRequest, err.Error())
}
