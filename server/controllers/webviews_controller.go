package controllers

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/commands/contact"
	"github.com/h3poteto/fascia/server/handlers"
	mailer "github.com/h3poteto/fascia/server/mailers/inquiry_mailer"
	"github.com/h3poteto/fascia/server/session"
	"github.com/h3poteto/fascia/server/validators"

	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
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

	return c.Render(http.StatusOK, "webviews/oauth_sign_in.html.tpl", map[string]interface{}{
		"title":      "SignIn",
		"privateURL": privateURL,
		"publicURL":  publicURL,
	})
}

// SignIn is a sign in action for mobile app
func (u *Webviews) SignIn(c echo.Context) error {
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

	return c.Render(http.StatusOK, "webviews/sign_in.html.tpl", map[string]interface{}{
		"title": "SignIn",
		"token": token,
	})
}

// NewSession is a sign in action for mobile app
func (u *Webviews) NewSession(c echo.Context) error {
	err := session.SharedInstance().Clear(c.Request(), c.Response())
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	var signInForm SignInForm
	err = c.Bind(signInForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	if !CheckCSRFToken(c, signInForm.Token) {
		err := errors.New("cannot verify CSRF token")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	userService, err := handlers.LoginUser(template.HTMLEscapeString(signInForm.Email), template.HTMLEscapeString(signInForm.Password))
	if err != nil {
		logging.SharedInstance().Controller(c).Infof("login error: %v", err)
		return c.Redirect(http.StatusFound, "/webviews/sign_in")
	}
	logging.SharedInstance().Controller(c).Debugf("login success: %+v", userService)

	option := &sessions.Options{Path: "/", MaxAge: config.Element("session").(map[interface{}]interface{})["timeout"].(int)}
	err = session.SharedInstance().Set(c.Request(), c.Response(), "current_user_id", userService.UserEntity.ID, option)
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("login success")
	return c.Redirect(http.StatusFound, "/webviews/callback")
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

	command := contact.InitCreateInquiry(0, newInquiryFrom.Email, newInquiryFrom.Name, newInquiryFrom.Message)
	if err := command.Run(); err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("success to create inquiry")

	// ここでemail送信
	go mailer.Notify(command.InquiryEntity)
	logging.SharedInstance().Controller(c).Info("success to send inquiry")
	return c.Render(http.StatusCreated, "webviews/inquiries/create.html.tpl", map[string]interface{}{
		"title": "Contact",
	})
}
