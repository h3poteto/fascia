package controllers

import (
	"../mailers/password_mailer"
	"../models/reset_password"
	"../models/user"
	"../modules/logging"
	"github.com/flosch/pongo2"
	"github.com/goji/param"
	"github.com/zenazn/goji/web"
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
		logging.SharedInstance().MethodInfo("PasswordsController", "New", true, c).Errorf("CSRF error: %v", err)
		InternalServerError(w, r)
		return
	}
	tpl, err := pongo2.DefaultSet.FromFile("new_password.html.tpl")
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "New", true, c).Errorf("template error: %v", err)
		InternalServerError(w, r)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "PasswordReset", "token": token}, w)
}

func (u *Passwords) Create(c web.C, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Create", true, c).Errorf("wrong form: %v", err)
		BadRequest(w, r)
		return
	}
	var newPasswordForm NewPasswordForm
	err = param.Parse(r.PostForm, &newPasswordForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Create", true, c).Errorf("wrong parameter: %v", err)
		InternalServerError(w, r)
		return
	}

	if !CheckCSRFToken(r, newPasswordForm.Token) {
		logging.SharedInstance().MethodInfo("PasswordsController", "Create", true, c).Error("cannot verify CSRF token")
		InternalServerError(w, r)
		return
	}

	targetUser, err := user.FindByEmail(newPasswordForm.Email)
	if err != nil {
		// OKにしておかないとEmail探りに使われる
		logging.SharedInstance().MethodInfo("PasswordsController", "Create", false, c).Infof("cannot find user: %v", err)
		http.Redirect(w, r, "/sign_in", 302)
		return
	}

	reset := reset_password.GenerateResetPassword(targetUser.ID, targetUser.Email)
	if err := reset.Save(); err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Create", true, c).Error("password_reset save error: %v", err)
		InternalServerError(w, r)
		return
	}
	// ここでemail送信
	go password_mailer.Reset(reset.ID, targetUser.Email, reset.Token)
	http.Redirect(w, r, "/sign_in", 302)
	logging.SharedInstance().MethodInfo("PasswordsController", "Create", false, c).Info("success to send password reset request")
	return
}

func (u *Passwords) Edit(c web.C, w http.ResponseWriter, r *http.Request) {
	token, err := GenerateCSRFToken(c, w, r)
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Edit", true, c).Errorf("CSRF error: %v", err)
		InternalServerError(w, r)
		return
	}
	resetToken := r.URL.Query().Get("token")
	id, err := strconv.ParseInt(c.URLParams["id"], 10, 64)
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Edit", true, c).Errorf("parse error: %v", err)
		http.Error(w, "reset password not found", 404)
		return
	}
	if err := reset_password.Authenticate(id, resetToken); err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Edit", false, c).Info("cannot authenticate reset password: %v", err)
		InternalServerError(w, r)
		return
	}
	tpl, err := pongo2.DefaultSet.FromFile("edit_password.html.tpl")
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Edit", true, c).Errorf("template error: %v", err)
		InternalServerError(w, r)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "PasswordReset", "token": token, "id": id, "resetToken": resetToken}, w)
}

func (u *Passwords) Update(c web.C, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Update", true, c).Errorf("wrong form: %v", err)
		BadRequest(w, r)
		return
	}
	var editPasswordForm EditPasswordForm
	err = param.Parse(r.PostForm, &editPasswordForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Update", true, c).Errorf("wrong parameters: %v", err)
		InternalServerError(w, r)
		return
	}

	if !CheckCSRFToken(r, editPasswordForm.Token) {
		logging.SharedInstance().MethodInfo("PasswordsController", "Update", true, c).Errorf("cannot verify CSRF token")
		InternalServerError(w, r)
		return
	}

	if editPasswordForm.Password != editPasswordForm.PasswordConfirm {
		logging.SharedInstance().MethodInfo("PasswordsController", "Update", true, c).Error("cannot confirm password")
		InternalServerError(w, r)
		return
	}

	id, err := strconv.ParseInt(c.URLParams["id"], 10, 64)
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Update", true, c).Errorf("parse error: %v", err)
		http.Error(w, "reset password not found", 404)
		return
	}
	targetUser, err := reset_password.ChangeUserPassword(id, editPasswordForm.ResetToken, editPasswordForm.Password)
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Update", false, c).Info("cannot authenticate reset password")
		InternalServerError(w, r)
		return
	}

	go password_mailer.Changed(targetUser.Email)
	logging.SharedInstance().MethodInfo("PasswordsController", "Update", false, c).Info("success to change password")
	http.Redirect(w, r, "/sign_in", 302)
	return
}
