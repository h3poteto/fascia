package controllers

import (
	"../config"
	userModel "../models/user"
	"../modules/logging"

	"html/template"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/goji/param"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
)

type Sessions struct {
}

type SignInForm struct {
	Email    string `param:"email"`
	Password string `param:"password"`
	Token    string `param:"token"`
}

func (u *Sessions) SignIn(c web.C, w http.ResponseWriter, r *http.Request) {
	url := githubOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	token, err := GenerateCSRFToken(c, w, r)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("SessionsController", "SignIn", err, c).Errorf("CSRF error: %v", err)
		InternalServerError(w, r)
		return
	}
	tpl, err := pongo2.DefaultSet.FromFile("sign_in.html.tpl")
	if err != nil {
		err := errors.Wrap(err, "template error")
		logging.SharedInstance().MethodInfoWithStacktrace("SessionsController", "SignIn", err, c).Error(err)
		InternalServerError(w, r)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "SignIn", "oauthURL": url, "token": token}, w)
}

func (u *Sessions) NewSession(c web.C, w http.ResponseWriter, r *http.Request) {
	// 旧セッションの削除
	session, err := cookieStore.Get(r, "fascia")
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().MethodInfoWithStacktrace("SessionsController", "NewSession", err, c).Error(err)
		InternalServerError(w, r)
		return
	}
	session.Options = &sessions.Options{MaxAge: -1}
	err = session.Save(r, w)
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().MethodInfoWithStacktrace("SessionsController", "NewSession", err, c).Error(err)
		InternalServerError(w, r)
		return
	}
	err = r.ParseForm()
	if err != nil {
		err := errors.Wrap(err, "wrong form")
		logging.SharedInstance().MethodInfoWithStacktrace("SessionsController", "NewSession", err, c).Error(err)
		BadRequest(w, r)
		return
	}
	var signInForm SignInForm
	err = param.Parse(r.PostForm, &signInForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("SessionsController", "NewSession", err, c).Error(err)
		InternalServerError(w, r)
		return
	}

	if !CheckCSRFToken(r, signInForm.Token) {
		err := errors.New("cannot verify CSRF token")
		logging.SharedInstance().MethodInfoWithStacktrace("SessionsController", "NewSession", err, c).Error(err)
		InternalServerError(w, r)
		return
	}

	currentUser, err := userModel.Login(template.HTMLEscapeString(signInForm.Email), template.HTMLEscapeString(signInForm.Password))
	if err != nil {
		logging.SharedInstance().MethodInfo("SessionsController", "NewSession", c).Infof("login error: %v", err)
		http.Redirect(w, r, "/sign_in", 302)
		return
	}
	logging.SharedInstance().MethodInfo("SessionsController", "NewSession", c).Debugf("login success: %+v", currentUser)
	session, err = cookieStore.Get(r, "fascia")
	session.Options = &sessions.Options{Path: "/", MaxAge: config.Element("session").(map[interface{}]interface{})["timeout"].(int)}
	session.Values["current_user_id"] = currentUser.ID
	err = session.Save(r, w)
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().MethodInfoWithStacktrace("SessionsController", "NewSessions", err, c).Error(err)
		InternalServerError(w, r)
		return
	}
	logging.SharedInstance().MethodInfo("SessionsController", "NewSession", c).Info("login success")
	http.Redirect(w, r, "/", 302)
	return
}

func (u *Sessions) SignOut(c web.C, w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, "fascia")
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().MethodInfoWithStacktrace("SessionsController", "SignOut", err, c).Error(err)
		InternalServerError(w, r)
		return
	}
	session.Options = &sessions.Options{MaxAge: -1}
	err = session.Save(r, w)
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().MethodInfoWithStacktrace("SessionsController", "SignOut", err, c).Error(err)
		InternalServerError(w, r)
		return
	}
	logging.SharedInstance().MethodInfo("SessionsController", "SignOut", c).Info("logout success")
	http.Redirect(w, r, "/sign_in", 302)
	return
}

func (u *Sessions) Update(c web.C, w http.ResponseWriter, r *http.Request) {
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("SessionsController", "Update", c).Infof("login error: %v", err)
		http.Error(w, "Authentication Error", 401)
		return
	}
	logging.SharedInstance().MethodInfo("SessionsController", "Update", c).Info("login success")
	session, err := cookieStore.Get(r, "fascia")
	session.Options = &sessions.Options{
		Path:   "/",
		MaxAge: config.Element("session").(map[interface{}]interface{})["timeout"].(int),
	}
	session.Values["current_user_id"] = currentUser.ID
	err = session.Save(r, w)
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().MethodInfoWithStacktrace("SessionsController", "Update", err, c).Error(err)
		InternalServerError(w, r)
		return
	}
	logging.SharedInstance().MethodInfo("SessionsController", "Update", c).Info("session update success")
	return
}
