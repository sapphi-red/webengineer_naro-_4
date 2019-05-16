package util

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

type ResponseData struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func Return400(c echo.Context, err error) error {
	return c.String(http.StatusBadRequest, err.Error())
}

func Return500(c echo.Context, err error) error {
	return c.String(http.StatusInternalServerError, err.Error())
}

func ReturnDBError(c echo.Context, err error) error {
	return Return500(c, fmt.Errorf("db error: %v", err))
}

func ReturnErrorJSON(c echo.Context, content string) error {
	return c.JSON(http.StatusBadRequest, ResponseData{
		Type:    "Error",
		Content: content,
	})
}

func ReturnSuccessJSON(c echo.Context) error {
	return c.JSON(http.StatusOK, ResponseData{
		Type:    "Success",
		Content: "Success",
	})
}
