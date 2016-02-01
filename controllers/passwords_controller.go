package controllers

import (
	"../mailers/password_mailer"
	"../models/reset_password"
	"../models/user"
	"../modules/logging"
	"github.com/flosch/pongo2"
	"github.com/goji/param"
	"github.com/zenazn/goji/web"
	//"html/template"
	"net/http"
	"strconv"
)

type Passwords struct {
}

type NewPasswordForm struct {
	Email string `param:"email"`
	Token string `param:"token"`
}

type EditPasswordForm struct {
	Token           string `param:"token"`
	ResetToken      string `param:"reset-token"`
	Password        string `param:"password"`
	PasswordConfirm string `param:"password-confirm"`
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
	// ここでemail送信
	go password_mailer.Reset(targetUser.Email, reset.Token)
	http.Redirect(w, r, "/sign_in", 302)
	logging.SharedInstance().MethodInfo("PasswordsController", "Create").Info("success to send password reset request")
	return
}

func (u *Passwords) Edit(c web.C, w http.ResponseWriter, r *http.Request) {
	token, err := GenerateCSRFToken(c, w, r)
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Edit").Errorf("CSRF error: %v", err)
		http.Error(w, "CSRF error", 500)
		return
	}
	resetToken := r.URL.Query().Get("token")
	id, _ := strconv.ParseInt(c.URLParams["id"], 10, 64)
	if !reset_password.Authenticate(id, resetToken) {
		logging.SharedInstance().MethodInfo("PasswordsController", "Edit").Info("cannot authenticate reset password")
		http.Error(w, "token error", 500)
		return
	}
	tpl, err := pongo2.DefaultSet.FromFile("edit_password.html.tpl")
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Edit").Errorf("template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "PasswordReset", "token": token, "id": id, "resetToken": resetToken}, w)
}

func (u *Passwords) Update(c web.C, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Update").Errorf("wrong form: %v", err)
		http.Error(w, "Wrong Form", http.StatusInternalServerError)
		return
	}
	var editPasswordForm EditPasswordForm
	err = param.Parse(r.PostForm, &editPasswordForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Update").Errorf("wrong parameters: %v", err)
		http.Error(w, "Wrong Parameter", http.StatusInternalServerError)
		return
	}

	if !CheckCSRFToken(r, editPasswordForm.Token) {
		logging.SharedInstance().MethodInfo("PasswordsController", "Update").Errorf("cannot verify CSRF token")
		http.Error(w, "Cannot verify CSRF token", http.StatusInternalServerError)
		return
	}

	if editPasswordForm.Password != editPasswordForm.PasswordConfirm {
		logging.SharedInstance().MethodInfo("PasswordsController", "Update").Error("cannot confirm password")
		http.Error(w, "Password is invalid", http.StatusInternalServerError)
		return
	}

	id, _ := strconv.ParseInt(c.URLParams["id"], 10, 64)
	targetUser, err := reset_password.ChangeUserPassword(id, editPasswordForm.ResetToken, editPasswordForm.Password)
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Update").Info("cannot authenticate reset password")
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}

	go password_mailer.Changed(targetUser.Email)
	logging.SharedInstance().MethodInfo("PasswordsController", "Update").Info("success to change password")
	http.Redirect(w, r, "/sign_in", 302)
	return
}
