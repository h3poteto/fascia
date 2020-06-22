package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/middlewares"
	"github.com/h3poteto/fascia/server/session"
	"github.com/h3poteto/fascia/server/usecases/account"
	"github.com/h3poteto/fascia/server/views"

	"net/http"
	"net/url"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// RecaptchaSiteverifyURL is an url for google reCAPTCHA v3.
var RecaptchaSiteverifyURL = "https://www.google.com/recaptcha/api/siteverify"

// Sessions is controller struct for sessions
type Sessions struct {
}

// NewSessionForm is form object.
type NewSessionForm struct {
	Email             string `json:"email" form:"email"`
	Password          string `json:"password" form:"password"`
	RecaptchaResponse string `json:"recaptcha_response" form:"recaptcha_response"`
}

// RecaptchaResult contains response from recaptcha stieverify.
// failed: "{\n  \"success\": false,\n  \"error-codes\": [\n    \"missing-input-response\"\n  ]\n}"
// success: "{\n  \"success\": true,\n  \"challenge_ts\": \"2020-06-16T14:17:18Z\",\n  \"hostname\": \"localhost\",\n  \"score\": 0.9,\n  \"action\": \"contact\"\n
type RecaptchaResult struct {
	Success bool `json:"success"`
	*RecaptchaFailed
	*RecaptchaSuccess
}

// RecaptchaFailed defines failed response.
type RecaptchaFailed struct {
	ErrorCodes []string `json:"error-codes"`
}

// RecaptchaSuccess defines succeeded response.
type RecaptchaSuccess struct {
	Hostname string  `json:"hostname"`
	Score    float64 `json:"score"`
	Action   string  `json:"action"`
}

// SignIn renders a sign in form
func (u *Sessions) SignIn(c echo.Context) error {
	token, err := GenerateCSRFToken(c)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Errorf("CSRF error: %v", err)
		return err
	}

	return c.Render(http.StatusOK, "sign_in.html.tpl", map[string]interface{}{
		"title": "SignIn",
		"token": token,
	})
}

// Create sign in a new user.
func (u *Sessions) Create(c echo.Context) error {
	newSessionForm := new(NewSessionForm)
	if err := c.Bind(newSessionForm); err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	val := url.Values{}
	val.Add("secret", os.Getenv("RECAPTCHA_SECRET_KEY"))
	val.Add("response", newSessionForm.RecaptchaResponse)
	resp, err := http.PostForm(RecaptchaSiteverifyURL, val)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	jsonBytes := ([]byte)(byteArray)
	result := new(RecaptchaResult)
	if err := json.Unmarshal(jsonBytes, result); err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	if !result.Success || result.RecaptchaSuccess == nil {
		err := fmt.Errorf("Recaptcha failed: %v", result.RecaptchaFailed.ErrorCodes)
		logging.SharedInstance().ControllerWithStacktrace(err, c).Warn(err)
		return c.Redirect(http.StatusFound, "/sign_in")
	}
	if result.RecaptchaSuccess != nil && result.RecaptchaSuccess.Score < 0.5 {
		err := errors.New("Recaptcha score is too low")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Warn(err)
		return c.Redirect(http.StatusFound, "/sign_in")
	}

	user, err := account.Authenticate(newSessionForm.Email, newSessionForm.Password)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Warn(err)
		return c.Redirect(http.StatusFound, "/sign_in")
	}
	option := &sessions.Options{
		Path:     "/",
		MaxAge:   config.Element("session").(map[interface{}]interface{})["timeout"].(int),
		HttpOnly: true,
	}
	if err := session.SharedInstance().Set(c.Request(), c.Response(), "current_user_id", user.ID, option); err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return c.Redirect(http.StatusFound, "/sign_in")
	}
	return c.Redirect(http.StatusFound, "/")
}

// SignOut delete a session and logout
func (u *Sessions) SignOut(c echo.Context) error {
	err := session.SharedInstance().Clear(c.Request(), c.Response())
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	logging.SharedInstance().Controller(c).Info("logout success")
	return c.JSON(http.StatusOK, nil)
}

// Update a session
func (u *Sessions) Update(c echo.Context) error {
	uc, ok := c.(*middlewares.LoginContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	userService := uc.CurrentUser

	option := &sessions.Options{
		Path:     "/",
		MaxAge:   config.Element("session").(map[interface{}]interface{})["timeout"].(int),
		HttpOnly: true,
	}
	err := session.SharedInstance().Set(c.Request(), c.Response(), "current_user_id", userService.ID, option)
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("session update success")
	return c.JSON(http.StatusOK, nil)
}

// Show returns current session and login user
func (u *Sessions) Show(c echo.Context) error {
	uc, ok := c.(*middlewares.LoginContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	userEntity := uc.CurrentUser
	jsonUser, err := views.ParseUserJSON(userEntity)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	return c.JSON(http.StatusOK, jsonUser)
}
