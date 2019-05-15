package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/srinathgs/mysqlstore"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type LoginRequestBody struct {
	Username string `json:"username,omitempty"  db:"username"`
	Password string `json:"password,omitempty"  db:"password"`
}

type User struct {
	Username   string `json:"username,omitempty"  db:"Username"`
	HashedPass string `json:"-"  db:"HashedPass"`
}

type Country struct {
	Code           string  `json:"code,omitempty"  db:"Code"`
	Name           string  `json:"name,omitempty"  db:"Name"`
	Continent      string  `json:"continent,omitempty"  db:"Continent"`
	Region         string  `json:"region,omitempty"  db:"Region"`
	SurfaceArea    float64 `json:"surface_area,omitempty"  db:"SurfaceArea"`
	IndepYear      int     `json:"indep_year,omitempty"  db:"IndepYear"`
	Population     int     `json:"population,omitempty"  db:"Population"`
	LifeExpectancy float64 `json:"life_expectancy,omitempty"  db:"LifeExpectancy"`
	GNP            float64 `json:"GNP,omitempty"  db:"GNP"`
	GNPOld         float64 `json:"GNP_old,omitempty"  db:"GNPOld"`
	LocalName      string  `json:"local_name,omitempty"  db:"LocalName"`
	GovernmentForm string  `json:"government_form,omitempty"  db:"GovernmentForm"`
	HeadOfState    string  `json:"head_of_state,omitempty"  db:"HeadOfState"`
	Capital        int     `json:"capital,omitempty"  db:"Capital"`
	Code2          string  `json:"code2,omitempty"  db:"Code2"`
}

type City struct {
	ID          int    `json:"ID,omitempty" db:"ID"`
	Name        string `json:"name,omitempty" db:"Name"`
	CountryCode string `json:"country_code,omitempty" db:"CountryCode"`
	District    string `json:"district,omitempty" db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

type ResponseData struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

var (
	db *sqlx.DB
)

func main() {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOSTNAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	db = _db

	// ----
	// 認証
	// ----
	store, err := mysqlstore.NewMySQLStoreFromConnection(
		db.DB,
		"sessions",
		"/",
		60*60*24*14,
		[]byte("secret-token"),
	)
	if err != nil {
		panic(err)
	}
	// ----

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.POST("/signup", postSignUpHandler)
	e.POST("/login", postLoginHandler)

	withLogin := e.Group("")
	withLogin.Use(checkLogin)

	withLogin.GET("/cities/:cityName", getCityInfoHandler)
	withLogin.POST("/cities", postCityInfoHandler)
	withLogin.DELETE("/cities/:cityName", deleteCityInfoHandler)

	e.Start(":12100")
}

func postSignUpHandler(c echo.Context) error {
	req := LoginRequestBody{}
	c.Bind(&req)

	// Todo: more validation
	if req.Password == "" || req.Username == "" {
		// Todo: better error
		return c.String(http.StatusBadRequest, "項目が空です")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

	var count int
	err = db.Get(&count, `SELECT COUNT(*) FROM users WHERE Username = ?`, req.Username)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	if count > 0 {
		return c.String(http.StatusConflict, "ユーザーが既に存在しています")
	}

	_, err = db.Exec(
		`INSERT INTO users (Username, HashedPass) VALUES (?, ?)`,
		req.Username,
		hashedPass,
	)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	return c.NoContent(http.StatusCreated)
}

func postLoginHandler(c echo.Context) error {
	req := LoginRequestBody{}
	c.Bind(&req)

	user := User{}
	err := db.Get(
		&user,
		`SELECT * FROM users WHERE username = ?`,
		req.Username,
	)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(req.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return c.NoContent(http.StatusForbidden)
		} else {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}
	sess.Values["userName"] = req.Username
	sess.Save(c.Request(), c.Response())

	return c.NoContent(http.StatusOK)
}

func checkLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "something wrong in getting session")
		}

		if sess.Values["userName"] == nil {
			return c.String(http.StatusForbidden, "please login")
		}
		c.Set("userName", sess.Values["userName"].(string))

		return next(c)
	}
}

func getCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")

	city := City{}
	db.Get(
		&city,
		`SELECT * FROM city WHERE Name = ?`,
		cityName,
	)
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

	_, err = db.Exec(
		`INSERT INTO city (Name, CountryCode, District, Population) VALUES (?, ?, ?, ?)`,
		data.Name,
		data.CountryCode,
		data.District,
		data.Population,
	)
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

	_, err := db.Exec(
		`DELETE FROM city WHERE Name = ?`,
		cityName,
	)
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
