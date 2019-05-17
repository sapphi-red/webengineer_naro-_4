package router

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
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

const (
	USERNAME_MAX_LENGTH = 30
	PASSWORD_MAX_LENGTH = 72
	PEPPER              = "すごすごーい"
)

type loginHandler struct {
	db *sqlx.DB
}

// unneeded
func (h loginHandler) addSaltAndPepper(username, password string) string {
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

func (h loginHandler) validateInputs(req LoginRequestBody) error {
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

func (h loginHandler) SignUp(c echo.Context) error {
	req := LoginRequestBody{}
	c.Bind(&req)

	err := h.validateInputs(req)
	if err != nil {
		return return400(c, err)
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(h.addSaltAndPepper(req.Password, req.Username)), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

	var count int
	err = h.db.Get(&count, `SELECT COUNT(*) FROM users WHERE Username = ?`, req.Username)
	if err != nil {
		return returnDBError(c, "UserGettingError", err)
	}

	if count > 0 {
		return c.String(http.StatusConflict, "ユーザーが既に存在しています")
	}

	_, err = h.db.Exec(
		`INSERT INTO users (Username, HashedPass) VALUES (?, ?)`,
		req.Username,
		hashedPass,
	)
	if err != nil {
		return returnDBError(c, "UserAddingError", err)
	}

	return c.NoContent(http.StatusCreated)
}

func (h loginHandler) Login(c echo.Context) error {
	req := LoginRequestBody{}
	c.Bind(&req)

	user := User{}
	err := h.db.Get(
		&user,
		`SELECT * FROM users WHERE username = ?`,
		req.Username,
	)
	if err != nil {
		return returnDBError(c, "UserGettingError", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(h.addSaltAndPepper(req.Password, req.Username)))
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

func CreateLoginRoutes(e *echo.Echo, db *sqlx.DB) {
	handler := &loginHandler{db: db}
	e.POST("/signup", handler.SignUp)
	e.POST("/login", handler.Login)
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
