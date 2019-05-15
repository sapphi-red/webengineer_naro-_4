package main

import (
	"errors"
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

func validateInputs(req LoginRequestBody) error {
	if req.Username == "" {
		return errors.New("ユーザー名が空です")
	}
	if len(req.Username) > 30 {
		return errors.New("ユーザー名が長すぎます")
	}

	if req.Password == "" {
		return errors.New("パスワードが空です")
	}
	if len(req.Password) > 72 {
		return errors.New("パスワードが長すぎます")
	}
	return nil
}

func postSignUpHandler(c echo.Context) error {
	req := LoginRequestBody{}
	c.Bind(&req)

	err := validateInputs(req)
	if err != nil {
		return return400(c, err)
	}
	
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

	var count int
	err = db.Get(&count, `SELECT COUNT(*) FROM users WHERE Username = ?`, req.Username)
	if err != nil {
		return returnDBError(c, err)
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
		return returnDBError(c, err)
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
		return returnDBError(c, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(req.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return c.NoContent(http.StatusForbidden)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return return500(c, errors.New("something wrong in getting session"))
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
			return return500(c, errors.New("something wrong in getting session"))
		}

		if sess.Values["userName"] == nil {
			return c.String(http.StatusForbidden, "please login")
		}
		c.Set("userName", sess.Values["userName"].(string))

		return next(c)
	}
}
