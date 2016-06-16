package controllers

import (
	userModel "../models/user"
	"../modules/logging"
	"../validators"

	"html/template"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/goji/param"
	"github.com/pkg/errors"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
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
		logging.SharedInstance().MethodInfoWithStacktrace("RegistrationsController", "SignUp", err, c).Errorf("CSRF error: %v", err)
		InternalServerError(w, r)
		return
	}

	tpl, err := pongo2.DefaultSet.FromFile("sign_up.html.tpl")
	if err != nil {
		err := errors.Wrap(err, "template error")
		logging.SharedInstance().MethodInfoWithStacktrace("RegistrationsController", "SignUp", err, c).Error(err)
		InternalServerError(w, r)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "SignUp", "oauthURL": url, "token": token}, w)
}

func (u *Registrations) Registration(c web.C, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		err := errors.Wrap(err, "wrong form")
		logging.SharedInstance().MethodInfoWithStacktrace("RegistrationsController", "SignUp", err, c).Error(err)
		BadRequest(w, r)
		return
	}

	var signUpForm SignUpForm
	err = param.Parse(r.PostForm, &signUpForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("RegistrationsController", "SignUp", err, c).Error(err)
		InternalServerError(w, r)
		return
	}
	logging.SharedInstance().MethodInfo("RegistrationsController", "SignUp", c).Debugf("post registration form: %+v", signUpForm)

	if !CheckCSRFToken(r, signUpForm.Token) {
		err := errors.New("cannot verify CSRF token")
		logging.SharedInstance().MethodInfoWithStacktrace("RegistrationsController", "SignUp", err, c).Error(err)
		InternalServerError(w, r)
		return
	}

	// sign up
	valid, err := validators.UserRegistrationValidation(signUpForm.Email, signUpForm.Password, signUpForm.PasswordConfirm)
	// TODO: 失敗していることは何かしらの方法で伝えたい
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("RegistrationsController", "SignUp", c).Infof("validation failed: %v", err)
		http.Redirect(w, r, "/sign_up", 302)
		return
	}
	// TODO: ここCSRFのmiddlewareとかでなんとかならんかなぁ
	_, err = userModel.Registration(
		template.HTMLEscapeString(signUpForm.Email),
		template.HTMLEscapeString(signUpForm.Password),
		template.HTMLEscapeString(signUpForm.PasswordConfirm),
	)
	if err != nil {
		// TODO: 登録情報が間違っていることを通知したい
		logging.SharedInstance().MethodInfo("RegistrationsController", "SignUp", c).Infof("registration error: %v", err)
		http.Redirect(w, r, "/sign_up", 302)
		return
	}

	// TODO: 成功していることも伝えたい
	logging.SharedInstance().MethodInfo("RegistrationsController", "SignUp", c).Info("registration success")
	http.Redirect(w, r, "/sign_in", 302)
	return
}
