package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"net/http"

	"github.com/sapphi-red/webengineer_naro-_4/database"
	"github.com/sapphi-red/webengineer_naro-_4/login"
	"github.com/sapphi-red/webengineer_naro-_4/router"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := database.ConnectDB()
	store := login.CreateSessionStore(db)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	login.CreateLoginRouter(e, db)

	withLogin := e.Group("")
	withLogin.Use(login.CheckLogin)
	router.CreateRoutes(withLogin)

	e.Start(":12100")
}
