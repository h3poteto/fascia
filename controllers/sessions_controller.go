package controllers

import (
	userModel "../models/user"
	"../modules/logging"
	"github.com/flosch/pongo2"
	"github.com/goji/param"
	"github.com/gorilla/sessions"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
	"net/http"
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
		logging.SharedInstance().MethodInfo("SessionsController", "SignIn").Errorf("CSRF error: %v", err.Error())
		http.Error(w, "CSRF error", 500)
		return
	}
	tpl, err := pongo2.DefaultSet.FromFile("sign_in.html.tpl")
	if err != nil {
		logging.SharedInstance().MethodInfo("SessionsController", "SignIn").Errorf("template error: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "SignIn", "oauthURL": url, "token": token}, w)
}

func (u *Sessions) NewSession(c web.C, w http.ResponseWriter, r *http.Request) {
	// 旧セッションの削除
	session, err := cookieStore.Get(r, "fascia")
	session.Options = &sessions.Options{MaxAge: -1}
	session.Save(r, w)
	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("SessionsController", "NewSession").Errorf("wrong form: %v", err.Error())
		http.Error(w, "Wrong Form", 400)
		return
	}
	var signInForm SignInForm
	err = param.Parse(r.PostForm, &signInForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("SessionsController", "NewSession").Errorf("wrong parameter: %v", err.Error())
		http.Error(w, "Wrong Parameter", 500)
		return
	}

	if !CheckCSRFToken(r, signInForm.Token) {
		logging.SharedInstance().MethodInfo("SessionsController", "NewSession").Error("cannot verify CSRF token")
		http.Error(w, "Cannot verify CSRF token", 500)
		return
	}
	current_user, err := userModel.Login(signInForm.Email, signInForm.Password)
	if err != nil {
		logging.SharedInstance().MethodInfo("SessionsController", "NewSession").Errorf("login error: %v", err.Error())
		http.Redirect(w, r, "/sign_in", 302)
		return
	}
	logging.SharedInstance().MethodInfo("SessionsController", "NewSession").Debugf("login success: %+v", current_user)
	session, err = cookieStore.Get(r, "fascia")
	session.Options = &sessions.Options{Path: "/", MaxAge: 3600}
	session.Values["current_user_id"] = current_user.Id
	// TODO: err処理を増やす
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
	return
}

func (u *Sessions) SignOut(c web.C, w http.ResponseWriter, r *http.Request) {
	session, _ := cookieStore.Get(r, "fascia")
	session.Options = &sessions.Options{MaxAge: -1}
	// TODO: err処理を増やす
	session.Save(r, w)
	http.Redirect(w, r, "/sign_in", 302)
	return
}
