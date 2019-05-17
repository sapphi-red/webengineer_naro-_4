package router

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

type ResponseData struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func return400(c echo.Context, err error) error {
	return c.String(http.StatusBadRequest, err.Error())
}

func return500(c echo.Context, err error) error {
	return c.String(http.StatusInternalServerError, err.Error())
}

func returnDBError(c echo.Context, name string, err error) error {
	fmt.Printf(name+": %v\n", err)
	return return500(c, fmt.Errorf("エラーが発生しました"))
}

func returnErrorJSON(c echo.Context, content string) error {
	return c.JSON(http.StatusBadRequest, ResponseData{
		Type:    "Error",
		Content: content,
	})
}

func returnSuccessJSON(c echo.Context) error {
	return c.JSON(http.StatusOK, ResponseData{
		Type:    "Success",
		Content: "Success",
	})
}
