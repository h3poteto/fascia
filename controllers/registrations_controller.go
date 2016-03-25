package controllers

import (
	userModel "../models/user"
	"../modules/logging"
	"github.com/flosch/pongo2"
	"github.com/goji/param"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
	"html/template"
	"net/http"
)

type Registrations struct {
}

type SignUpForm struct {
	Email           string `param:"email"`
	Password        string `param:"password"`
	PasswordConfirm string `param:"password-confirm"`
	Token           string `param:"token"`
}

func (u *Registrations) SignUp(c web.C, w http.ResponseWriter, r *http.Request) {
	url := githubOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	token, err := GenerateCSRFToken(c, w, r)
	if err != nil {
		logging.SharedInstance().MethodInfo("RegistrationsController", "SignUp", true, c).Errorf("CSRF error: %v", err)
		InternalServerError(w, r)
		return
	}

	tpl, err := pongo2.DefaultSet.FromFile("sign_up.html.tpl")
	if err != nil {
		logging.SharedInstance().MethodInfo("RegistrationsController", "SignUp", true, c).Errorf("template error: %v", err)
		InternalServerError(w, r)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "SignUp", "oauthURL": url, "token": token}, w)
}

func (u *Registrations) Registration(c web.C, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("RegistrationsController", "SignUp", true, c).Errorf("wrong form: %v", err)
		BadRequest(w, r)
		return
	}

	var signUpForm SignUpForm
	err = param.Parse(r.PostForm, &signUpForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("RegistrationsController", "SignUp", true, c).Errorf("wrong parameter: %v", err)
		InternalServerError(w, r)
		return
	}
	logging.SharedInstance().MethodInfo("RegistrationsController", "SignUp", false, c).Debugf("post registration form: %+v", signUpForm)

	if !CheckCSRFToken(r, signUpForm.Token) {
		logging.SharedInstance().MethodInfo("RegistrationsController", "SignUp", true, c).Error("cannot verify CSRF token")
		InternalServerError(w, r)
		return
	}

	if signUpForm.Password == signUpForm.PasswordConfirm {
		// login
		_, err := userModel.Registration(template.HTMLEscapeString(signUpForm.Email), template.HTMLEscapeString(signUpForm.Password))
		if err != nil {
			// TODO: 二重登録ができないことを通知したい
			logging.SharedInstance().MethodInfo("RegistrationsController", "SignUp", true, c).Errorf("registration error: %v", err)
			http.Redirect(w, r, "/sign_up", 302)
			return
		} else {
			logging.SharedInstance().MethodInfo("RegistrationsController", "SignUp", false, c).Info("registration success")
			http.Redirect(w, r, "/sign_in", 302)
			return
		}
	} else {
		// error
		logging.SharedInstance().MethodInfo("RegistrationsController", "SignUp", true, c).Error("cannot confirm password")
		http.Redirect(w, r, "/sign_up", 302)
		return
	}
}
