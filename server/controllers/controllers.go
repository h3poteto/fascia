package controllers

import (
	"github.com/h3poteto/fascia/server/middlewares"
	"github.com/h3poteto/fascia/server/session"

	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var githubPrivateConf = &oauth2.Config{
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	Scopes:       []string{"repo", "write:repo_hook", "user:email"},
	Endpoint:     github.Endpoint,
}

var githubPublicConf = &oauth2.Config{
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	Scopes:       []string{"public_repo", "write:repo_hook", "user:email"},
	Endpoint:     github.Endpoint,
}

// CheckCSRFToken check token in session.
// To stub in test, I substitue the pointer of this function to global variable.
var CheckCSRFToken = checkCSRF

// GenerateCSRFToken prepare new CSRF token
var GenerateCSRFToken = generateCSRF

// NewJSONError prepare json error struct
var NewJSONError = middlewares.NewJSONError

// NewValidationError prepare validation error interface
var NewValidationError = middlewares.NewValidationError

// generateCSRF generate new CSRF token
func generateCSRF(c echo.Context) (string, error) {
	// Generate token from salt and current time.
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	io.WriteString(h, "secret_key_salt")
	token := fmt.Sprintf("%x", h.Sum(nil))

	err := session.SharedInstance().Set(c.Request(), c.Response(), "token", token)
	if err != nil {
		return "", errors.Wrap(err, "cookie error")
	}
	return token, nil
}

func checkCSRF(c echo.Context, token string) bool {
	t, err := session.SharedInstance().Get(c.Request(), "token")
	if err != nil {
		return false
	}

	if t.(string) != token {
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

// PrivacyPolicy render privacy policy html.
func PrivacyPolicy(c echo.Context) error {
	return c.Render(http.StatusOK, "privacy_policy.html.tpl", map[string]interface{}{
		"title": "PrivacyPolicy",
	})
}
