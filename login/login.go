package login

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/sapphi-red/webengineer_naro-_4/util"
	"github.com/srinathgs/mysqlstore"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type LoginRequestBody struct {
	Username string `json:"username,omitempty"  db:"username"`
	Password string `json:"password,omitempty"  db:"password"`
}

type User struct {
	Username   string `json:"username,omitempty"  db:"Username"`
	HashedPass string `json:"-"  db:"HashedPass"`
}

const USERNAME_MAX_LENGTH = 30
const PASSWORD_MAX_LENGTH = 72
const PEPPER = "すごすごーい"

func CreateSessionStore(db *sqlx.DB) *mysqlstore.MySQLStore {
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

func CreateLoginRouter(e *echo.Echo, db *sqlx.DB) {
	e.POST("/signup", makePostSignUpHandler(db))
	e.POST("/login", makePostLoginHandler(db))
}

func addSaltAndPepper(username, password string) string {
	if len(password) > PASSWORD_MAX_LENGTH {
		return password
	}

	newPassword := username + password
	if len(newPassword) > PASSWORD_MAX_LENGTH {
		return newPassword[:PASSWORD_MAX_LENGTH]
	}
	newPassword = newPassword + PEPPER
	if len(newPassword) > PASSWORD_MAX_LENGTH {
		return newPassword[:PASSWORD_MAX_LENGTH]
	}
	return newPassword
}

func validateInputs(req LoginRequestBody) error {
	if req.Username == "" {
		return errors.New("ユーザー名が空です")
	}
	if len(req.Username) > USERNAME_MAX_LENGTH {
		return errors.New("ユーザー名が長すぎます")
	}

	if req.Password == "" {
		return errors.New("パスワードが空です")
	}
	if len(req.Password) > PASSWORD_MAX_LENGTH {
		return errors.New("パスワードが長すぎます")
	}
	return nil
}

func makePostSignUpHandler(db *sqlx.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		req := LoginRequestBody{}
		c.Bind(&req)

		err := validateInputs(req)
		if err != nil {
			return util.Return400(c, err)
		}

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(addSaltAndPepper(req.Password, req.Username)), bcrypt.DefaultCost)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
		}

		var count int
		err = db.Get(&count, `SELECT COUNT(*) FROM users WHERE Username = ?`, req.Username)
		if err != nil {
			return util.ReturnDBError(c, err)
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
			return util.ReturnDBError(c, err)
		}

		return c.NoContent(http.StatusCreated)
	}
}

func makePostLoginHandler(db *sqlx.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		req := LoginRequestBody{}
		c.Bind(&req)

		user := User{}
		err := db.Get(
			&user,
			`SELECT * FROM users WHERE username = ?`,
			req.Username,
		)
		if err != nil {
			return util.ReturnDBError(c, err)
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(addSaltAndPepper(req.Password, req.Username)))
		if err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				return c.NoContent(http.StatusForbidden)
			}
			return c.NoContent(http.StatusInternalServerError)
		}

		sess, err := session.Get("sessions", c)
		if err != nil {
			fmt.Println(err)
			return util.Return500(c, errors.New("something wrong in getting session"))
		}
		sess.Values["userName"] = req.Username
		sess.Save(c.Request(), c.Response())

		return c.NoContent(http.StatusOK)
	}
}

func CheckLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			fmt.Println(err)
			return util.Return500(c, errors.New("something wrong in getting session"))
		}

		if sess.Values["userName"] == nil {
			return c.String(http.StatusForbidden, "please login")
		}
		c.Set("userName", sess.Values["userName"].(string))

		return next(c)
	}
}
