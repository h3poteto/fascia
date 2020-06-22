package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/h3poteto/fascia/lib/modules/logging"
	mailer "github.com/h3poteto/fascia/server/mailers/inquiry_mailer"
	"github.com/h3poteto/fascia/server/usecases/contact"
	"github.com/h3poteto/fascia/server/validators"

	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Webviews controller struct
type Webviews struct {
}

// OauthSignIn render oauth login page
func (u *Webviews) OauthSignIn(c echo.Context) error {
	privateURL := githubPrivateConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	publicURL := githubPublicConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	token, err := GenerateCSRFToken(c)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Errorf("CSRF error: %v", err)
		return err
	}

	// Set cookie for iOS application when authentication callback
	cookie := http.Cookie{
		Path:    "/",
		Name:    "fascia-ios",
		Value:   "login-session",
		Expires: time.Now().AddDate(0, 0, 1),
	}
	c.SetCookie(&cookie)

	return c.Render(http.StatusOK, "webviews/oauth_sign_in.html.tpl", map[string]interface{}{
		"title":      "SignIn",
		"privateURL": privateURL,
		"publicURL":  publicURL,
		"token":      token,
	})
}

// Callback is an empty page for mobile application handling
func (u *Webviews) Callback(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

// NewInquiry is a contact form for mobile app
func (u *Webviews) NewInquiry(c echo.Context) error {
	return c.Render(http.StatusOK, "webviews/inquiries/new.html.tpl", map[string]interface{}{
		"title": "Contact",
	})
}

// Inquiry create a new inquiry from contact
func (u *Webviews) Inquiry(c echo.Context) error {
	newInquiryFrom := new(NewInquiryForm)
	if err := c.Bind(newInquiryFrom); err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Debugf("post new inquiry parameter: %+v", newInquiryFrom)

	valid, err := validators.InquiryCreateValidation(newInquiryFrom.Email, newInquiryFrom.Name, newInquiryFrom.Message)
	if err != nil || !valid {
		logging.SharedInstance().Controller(c).Infof("validation error: %v", err)
		return c.Render(http.StatusUnprocessableEntity, "webviews/inquiries/new.html.tpl", map[string]interface{}{
			"title": "Contact",
			"error": err.Error(),
		})
	}

	val := url.Values{}
	val.Add("secret", os.Getenv("RECAPTCHA_SECRET_KEY"))
	val.Add("response", newInquiryFrom.RecaptchaResponse)
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
		return c.Render(http.StatusUnprocessableEntity, "webviews/inquiries/new.html.tpl", map[string]interface{}{
			"title": "Contact",
			"error": "Session error",
		})
	}
	if result.RecaptchaSuccess != nil && result.RecaptchaSuccess.Score < 0.5 {
		err := errors.New("Recaptcha score is too low")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Warn(err)
		return c.Render(http.StatusUnprocessableEntity, "webviews/inquiries/new.html.tpl", map[string]interface{}{
			"title": "Contact",
			"error": "Session error",
		})
	}

	inquiry, err := contact.CreateInquiry(newInquiryFrom.Email, newInquiryFrom.Name, newInquiryFrom.Message)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("success to create inquiry")

	// Send email
	go mailer.Notify(inquiry)
	logging.SharedInstance().Controller(c).Info("success to send inquiry")
	return c.Render(http.StatusCreated, "webviews/inquiries/create.html.tpl", map[string]interface{}{
		"title": "Contact",
	})
}
