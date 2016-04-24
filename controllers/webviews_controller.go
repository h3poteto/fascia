package controllers

import (
	"../config"
	userModel "../models/user"
	"../modules/logging"
	"time"

	"github.com/flosch/pongo2"
	"github.com/goji/param"
	"github.com/gorilla/sessions"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
	"html/template"
	"net/http"
)

// Webviews controller struct
type Webviews struct {
}

// SignIn is a sign in action for mobile app
func (u *Webviews) SignIn(c web.C, w http.ResponseWriter, r *http.Request) {
	url := githubOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	token, err := GenerateCSRFToken(c, w, r)
	if err != nil {
		logging.SharedInstance().MethodInfo("WebviewsController", "SignIn", true, c).Errorf("CSRF error: %v", err)
		InternalServerError(w, r)
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
		logging.SharedInstance().MethodInfo("WebviewsController", "SignIn", true, c).Errorf("template error: %v", err)
		InternalServerError(w, r)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "SignIn", "oauthURL": url, "token": token}, w)
	return
}

// NewSession is a sign in action for mobile app
func (u *Webviews) NewSession(c web.C, w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, "fascia")
	if err != nil {
		logging.SharedInstance().MethodInfo("WebviewsController", "NewSession", true, c).Errorf("get session error: %v", err)
		InternalServerError(w, r)
		return
	}
	session.Options = &sessions.Options{MaxAge: -1}
	err = session.Save(r, w)
	if err != nil {
		logging.SharedInstance().MethodInfo("WebviewsController", "NewSession", true, c).Errorf("save session error: %v", err)
		InternalServerError(w, r)
		return
	}
	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("WebviewsController", "NewSession", true, c).Errorf("wrong form: %v", err)
		BadRequest(w, r)
		return
	}
	var signInForm SignInForm
	err = param.Parse(r.PostForm, &signInForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("WebviewsController", "NewSession", true, c).Errorf("wrong parameter: %v", err)
		InternalServerError(w, r)
		return
	}

	if !CheckCSRFToken(r, signInForm.Token) {
		logging.SharedInstance().MethodInfo("WebviewsController", "NewSession", true, c).Error("cannot verify CSRF token")
		InternalServerError(w, r)
		return
	}

	currentUser, err := userModel.Login(template.HTMLEscapeString(signInForm.Email), template.HTMLEscapeString(signInForm.Password))
	if err != nil {
		logging.SharedInstance().MethodInfo("WebviewsController", "NewSession", false, c).Infof("login error: %v", err)
		http.Redirect(w, r, "/webviews/sign_in", 302)
		return
	}
	logging.SharedInstance().MethodInfo("WebviewsController", "NewSession", false, c).Debugf("login success: %+v", currentUser)
	session, err = cookieStore.Get(r, "fascia")
	session.Options = &sessions.Options{Path: "/", MaxAge: config.Element("session").(map[interface{}]interface{})["timeout"].(int)}
	session.Values["current_user_id"] = currentUser.ID
	err = session.Save(r, w)
	if err != nil {
		logging.SharedInstance().MethodInfo("WebviewsController", "NewSessions", true, c).Errorf("session error: %v", err)
		InternalServerError(w, r)
		return
	}
	logging.SharedInstance().MethodInfo("WebviewsController", "NewSession", false, c).Info("login success")
	http.Redirect(w, r, "/webviews/callback", 302)
	return
}

// Callback is a empty page for mobile application handling
func (u *Webviews) Callback(c web.C, w http.ResponseWriter, r *http.Request) {
	return
}
