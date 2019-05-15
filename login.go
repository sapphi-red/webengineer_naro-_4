package main

import (
	"github.com/labstack/echo-contrib/session"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"github.com/labstack/echo"
	"github.com/srinathgs/mysqlstore"
)

type LoginRequestBody struct {
	Username string `json:"username,omitempty"  db:"username"`
	Password string `json:"password,omitempty"  db:"password"`
}

type User struct {
	Username   string `json:"username,omitempty"  db:"Username"`
	HashedPass string `json:"-"  db:"HashedPass"`
}

func createSessionStore() *mysqlstore.MySQLStore {
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
	return store
}

func createLoginRouter(e *echo.Echo) {
	e.POST("/signup", postSignUpHandler)
	e.POST("/login", postLoginHandler)
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
