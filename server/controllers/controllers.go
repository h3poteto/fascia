package controllers

import (
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/services"

	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// Key defines session's key
const Key = "fascia.io"

// JSONError is a struct for http error
type JSONError struct {
	Code    int    `json:code`
	Message string `json:message`
}

// NewJSONError render error json response and return error
func NewJSONError(err error, code int, c echo.Context) error {
	c.JSON(code, &JSONError{
		Code:    code,
		Message: http.StatusText(code),
	})
	return err
}

var githubOauthConf = &oauth2.Config{
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	Scopes:       []string{"repo", "write:repo_hook", "user:email"},
	Endpoint:     github.Endpoint,
}

var cookieStore = sessions.NewCookieStore([]byte(os.Getenv("SECRET")))

// ここテストでstubするために関数ポインタをグローバル変数に代入しておきます．もしインスタンスメソッドではない関数をstubする方法があれば，書き換えて構わない．
// CheckCSRFToken check token in session
var CheckCSRFToken = checkCSRF

// LoginRequired check login session
var LoginRequired = CheckLogin

// GenerateCSRFToken prepare new CSRF token
var GenerateCSRFToken = generateCSRF

// CheckLogin authenticate user
// If unauthorized, return 401
func CheckLogin(c echo.Context) (*services.User, error) {
	session, err := cookieStore.Get(c.Request(), Key)
	if err != nil {
		return nil, errors.New("cookie error")
	}
	id := session.Values["current_user_id"]
	if id == nil {
		return nil, errors.New("not logined")
	}
	currentUser, err := handlers.FindUser(id.(int64))
	if err != nil {
		return nil, err
	}
	return currentUser, nil
}

// generateCSRF generate new CSRF token
func generateCSRF(c echo.Context) (string, error) {
	session, err := cookieStore.Get(c.Request(), Key)
	if err != nil {
		return "", errors.Wrap(err, "cookie error")
	}

	// 現在時間とソルトからトークンを生成
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	io.WriteString(h, "secret_key_salt")
	token := fmt.Sprintf("%x", h.Sum(nil))
	session.Values["token"] = token

	err = cookieStore.Save(c.Request(), c.Response(), session)
	if err != nil {
		return "", errors.Wrap(err, "cookie error")
	}
	return token, nil
}

func checkCSRF(c echo.Context, token string) bool {
	session, err := cookieStore.Get(c.Request(), Key)
	if err != nil {
		return false
	}

	if session.Values["token"] != token {
		return false
	}
	return true
}

// BadRequest render 400
func BadRequest(c echo.Context) error {
	return c.Render(http.StatusBadRequest, "400.html.tpl", map[string]interface{}{
		"title": "BadRequest",
	})
}

// NotFound render 404
func NotFound(c echo.Context) error {
	return c.Render(http.StatusNotFound, "404.html.tpl", map[string]interface{}{
		"title": "NotFound",
	})
}

// InternalServerError render 500
func InternalServerError(c echo.Context) error {
	return c.Render(http.StatusInternalServerError, "500.html.tpl", map[string]interface{}{
		"title": "InternalServerError",
	})
}
