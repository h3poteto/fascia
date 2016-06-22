package controllers

import (
	"../config"
	userModel "../models/user"
	"../modules/logging"
	"time"

	"html/template"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/goji/param"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
)

// Webviews controller struct
type Webviews struct {
}

// SignIn is a sign in action for mobile app
func (u *Webviews) SignIn(c web.C, w http.ResponseWriter, r *http.Request) {
	url := githubOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	token, err := GenerateCSRFToken(c, w, r)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("WebviewsController", "SignIn", err, c).Errorf("CSRF error: %v", err)
		errorBody := "Cookie error, please clear your browser cookie."
		InternalServerError(w, r, &errorBody)
		return
	}

	// prepare cookie
	cookie := http.Cookie{
		Path:    "/",
		Name:    "fascia-ios",
		Value:   "login-session",
		Expires: time.Now().AddDate(0, 0, 1),
	}
	http.SetCookie(w, &cookie)

	tpl, err := pongo2.DefaultSet.FromFile("webviews/sign_in.html.tpl")
	if err != nil {
		err := errors.Wrap(err, "template error")
		logging.SharedInstance().MethodInfoWithStacktrace("WebviewsController", "SignIn", err, c).Error(err)
		InternalServerError(w, r, nil)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "SignIn", "oauthURL": url, "token": token}, w)
	return
}

// NewSession is a sign in action for mobile app
func (u *Webviews) NewSession(c web.C, w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, "fascia")
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().MethodInfoWithStacktrace("WebviewsController", "NewSession", err, c).Error(err)
		errorBody := "Cookie error, please clear your browser cookie."
		InternalServerError(w, r, &errorBody)
		return
	}
	session.Options = &sessions.Options{MaxAge: -1}
	err = session.Save(r, w)
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().MethodInfoWithStacktrace("WebviewsController", "NewSession", err, c).Error(err)
		errorBody := "Cookie error, please clear your browser cookie."
		InternalServerError(w, r, &errorBody)
		return
	}
	err = r.ParseForm()
	if err != nil {
		err := errors.Wrap(err, "wrong form")
		logging.SharedInstance().MethodInfoWithStacktrace("WebviewsController", "NewSession", err, c).Error(err)
		BadRequest(w, r)
		return
	}
	var signInForm SignInForm
	err = param.Parse(r.PostForm, &signInForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("WebviewsController", "NewSession", err, c).Error(err)
		InternalServerError(w, r, nil)
		return
	}

	if !CheckCSRFToken(r, signInForm.Token) {
		err := errors.New("cannot verify CSRF token")
		logging.SharedInstance().MethodInfoWithStacktrace("WebviewsController", "NewSession", err, c).Error(err)
		InternalServerError(w, r, nil)
		return
	}

	currentUser, err := userModel.Login(template.HTMLEscapeString(signInForm.Email), template.HTMLEscapeString(signInForm.Password))
	if err != nil {
		logging.SharedInstance().MethodInfo("WebviewsController", "NewSession", c).Infof("login error: %v", err)
		http.Redirect(w, r, "/webviews/sign_in", 302)
		return
	}
	logging.SharedInstance().MethodInfo("WebviewsController", "NewSession", c).Debugf("login success: %+v", currentUser)
	session, err = cookieStore.Get(r, "fascia")
	session.Options = &sessions.Options{Path: "/", MaxAge: config.Element("session").(map[interface{}]interface{})["timeout"].(int)}
	session.Values["current_user_id"] = currentUser.ID
	err = session.Save(r, w)
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().MethodInfoWithStacktrace("WebviewsController", "NewSessions", err, c).Error(err)
		errorBody := "Cookie error, please clear your browser cookie."
		InternalServerError(w, r, &errorBody)
		return
	}
	logging.SharedInstance().MethodInfo("WebviewsController", "NewSession", c).Info("login success")
	http.Redirect(w, r, "/webviews/callback", 302)
	return
}

// Callback is a empty page for mobile application handling
func (u *Webviews) Callback(c web.C, w http.ResponseWriter, r *http.Request) {
	return
}
