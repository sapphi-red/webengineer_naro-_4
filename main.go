package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"net/http"

	"github.com/sapphi-red/webengineer_naro-_4/database"
	"github.com/sapphi-red/webengineer_naro-_4/router"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := database.ConnectDB()
	store := database.CreateSessionStore()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	router.CreateLoginRoutes(e, db)
	router.CreateRoutes(e, db)

	e.Start(":12100")
}
