package controllers

import (
	"../models/reset_password"
	"../models/user"
	"../modules/logging"
	"github.com/flosch/pongo2"
	//"github.com/gorilla/sessions"
	"github.com/goji/param"
	"github.com/zenazn/goji/web"
	//"html/template"
	"net/http"
)

type Passwords struct {
}

type NewPasswordForm struct {
	Email string `param:"email"`
	Token string `param:"token"`
}

// tokenを発行し，expireと合わせてreset_passwordモデルにDB保存する
// idとtokenをメールで送る
// idとtoken, expireがあっていたらpasswordの編集を許可する
// passwordを新たに保存する
func (u *Passwords) New(c web.C, w http.ResponseWriter, r *http.Request) {
	token, err := GenerateCSRFToken(c, w, r)
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "New").Errorf("CSRF error: %v", err)
		http.Error(w, "CSRF error", 500)
		return
	}
	tpl, err := pongo2.DefaultSet.FromFile("new_password.html.tpl")
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "New").Errorf("template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "PasswordReset", "token": token}, w)
}

func (u *Passwords) Create(c web.C, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Create").Errorf("wrong form: %v", err)
		http.Error(w, "Wrong Form", 500)
		return
	}
	var newPasswordForm NewPasswordForm
	err = param.Parse(r.PostForm, &newPasswordForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Create").Errorf("wrong parameter: %v", err)
		http.Error(w, "Wrong Parameter", 500)
		return
	}

	if !CheckCSRFToken(r, newPasswordForm.Token) {
		logging.SharedInstance().MethodInfo("PasswordsController", "Create").Error("cannot verify CSRF token")
		http.Error(w, "Cannot verify CSRF token", 500)
		return
	}

	targetUser, err := user.FindByEmail(newPasswordForm.Email)
	if err != nil {
		// OKにしておかないとEmail探りに使われる
		logging.SharedInstance().MethodInfo("PasswordsController", "Create").Infof("cannot find user: %v", err)
		http.Redirect(w, r, "/sign_in", 302)
		return
	}

	reset := reset_password.GenerateResetPassword(targetUser.Id, targetUser.Email)
	if !reset.Save() {
		logging.SharedInstance().MethodInfo("PasswordsController", "Create").Error("password_reset save error")
		http.Error(w, "save error", 500)
		return
	}
	http.Redirect(w, r, "/sign_in", 302)
	return
}

func (u *Passwords) Edit(c web.C, w http.ResponseWriter, r *http.Request) {
}

func (u *Passwords) Update(c web.C, w http.ResponseWriter, r *http.Request) {
}
