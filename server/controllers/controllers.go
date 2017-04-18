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

	"github.com/ipfans/echo-session"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type JSONError struct {
	Code    int    `json:code`
	Message string `json:message`
}

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

// ここテストでstubするために関数ポインタをグローバル変数に代入しておきます．もしインスタンスメソッドではない関数をstubする方法があれば，書き換えて構わない．
var CheckCSRFToken = checkCSRF
var LoginRequired = CheckLogin

// CheckLogin authenticate user
// If unauthorized, return 401
func CheckLogin(c echo.Context) (*services.User, error) {
	session := session.Default(c)
	id := session.Get("current_user_id")
	if id == nil {
		return nil, errors.New("not logined")
	}
	currentUser, err := handlers.FindUser(id.(int64))
	if err != nil {
		return nil, err
	}
	return currentUser, nil
}

// GenerateCSRFToken generate new CSRF token
func GenerateCSRFToken(c echo.Context) (string, error) {
	session := session.Default(c)

	// 現在時間とソルトからトークンを生成
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	io.WriteString(h, "secret_key_salt")
	token := fmt.Sprintf("%x", h.Sum(nil))
	session.Set("token", token)

	err := session.Save()
	if err != nil {
		return "", errors.Wrap(err, "cookie error")
	}
	return token, nil
}

func checkCSRF(c echo.Context, token string) bool {
	session := session.Default(c)

	if session.Get("token") != token {
		return false
	}
	return true
}

func BadRequest(c echo.Context) error {
	return c.Render(http.StatusBadRequest, "400.html.tpl", map[string]interface{}{
		"title": "BadRequest",
	})
}

func NotFound(c echo.Context) error {
	return c.Render(http.StatusNotFound, "404.html.tpl", map[string]interface{}{
		"title": "NotFound",
	})
}

func InternalServerError(c echo.Context) error {
	return c.Render(http.StatusInternalServerError, "500.html.tpl", map[string]interface{}{
		"title": "InternalServerError",
	})
}
