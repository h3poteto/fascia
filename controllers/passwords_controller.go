package controllers

import (
	"../mailers/password_mailer"
	"../models/reset_password"
	"../models/user"
	"../modules/logging"
	"../validators"

	"net/http"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/goji/param"
	"github.com/pkg/errors"
	"github.com/zenazn/goji/web"
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
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "New", err, c).Error(err)
		InternalServerError(w, r)
		return
	}
	tpl, err := pongo2.DefaultSet.FromFile("new_password.html.tpl")
	if err != nil {
		err := errors.Wrap(err, "template error")
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "New", err, c).Error(err)
		InternalServerError(w, r)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "PasswordReset", "token": token}, w)
}

func (u *Passwords) Create(c web.C, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		err := errors.Wrap(err, "wrong form")
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Create", err, c).Error(err)
		BadRequest(w, r)
		return
	}
	var newPasswordForm NewPasswordForm
	err = param.Parse(r.PostForm, &newPasswordForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Create", err, c).Error(err)
		InternalServerError(w, r)
		return
	}

	if !CheckCSRFToken(r, newPasswordForm.Token) {
		err := errors.New("cannot verify CSRF token")
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Create", err, c).Error(err)
		InternalServerError(w, r)
		return
	}

	valid, err := validators.PasswordCreateValidation(newPasswordForm.Email)
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("PasswordsController", "Create", c).Infof("validation failed: %v", err)
		http.Redirect(w, r, "/passwords/new", 302)
		return
	}

	targetUser, err := user.FindByEmail(newPasswordForm.Email)
	if err != nil {
		// OKにしておかないとEmail探りに使われる
		logging.SharedInstance().MethodInfo("PasswordsController", "Create", c).Infof("cannot find user: %v", err)
		http.Redirect(w, r, "/sign_in", 302)
		return
	}

	reset := reset_password.GenerateResetPassword(targetUser.ID, targetUser.Email)
	if err := reset.Save(); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Create", err, c).Error(err)
		InternalServerError(w, r)
		return
	}
	// ここでemail送信
	go password_mailer.Reset(reset.ID, targetUser.Email, reset.Token)
	http.Redirect(w, r, "/sign_in", 302)
	logging.SharedInstance().MethodInfo("PasswordsController", "Create", c).Info("success to send password reset request")
	return
}

func (u *Passwords) Edit(c web.C, w http.ResponseWriter, r *http.Request) {
	token, err := GenerateCSRFToken(c, w, r)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Edit", err, c).Error(err)
		InternalServerError(w, r)
		return
	}
	resetToken := r.URL.Query().Get("token")
	id, err := strconv.ParseInt(c.URLParams["id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Edit", err, c).Error(err)
		http.Error(w, "reset password not found", 404)
		return
	}
	if err := reset_password.Authenticate(id, resetToken); err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Edit", c).Info("cannot authenticate reset password: %v", err)
		InternalServerError(w, r)
		return
	}
	tpl, err := pongo2.DefaultSet.FromFile("edit_password.html.tpl")
	if err != nil {
		err := errors.Wrap(err, "template error")
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Edit", err, c).Error(err)
		InternalServerError(w, r)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "PasswordReset", "token": token, "id": id, "resetToken": resetToken}, w)
}

func (u *Passwords) Update(c web.C, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		err := errors.Wrap(err, "wrong form")
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Update", err, c).Error(err)
		BadRequest(w, r)
		return
	}
	var editPasswordForm EditPasswordForm
	err = param.Parse(r.PostForm, &editPasswordForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameters")
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Update", err, c).Error(err)
		InternalServerError(w, r)
		return
	}

	if !CheckCSRFToken(r, editPasswordForm.Token) {
		err := errors.New("cannot verify CSRF token")
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Update", err, c).Error(err)
		InternalServerError(w, r)
		return
	}

	id, err := strconv.ParseInt(c.URLParams["id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Update", err, c).Error(err)
		http.Error(w, "reset password not found", 404)
		return
	}

	valid, err := validators.PasswordUpdateValidation(editPasswordForm.ResetToken, editPasswordForm.Password, editPasswordForm.PasswordConfirm)
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("PasswordController", "Update", c).Infof("validation failed: %v", err)
		http.Redirect(w, r, "/passwords/"+string(id)+"/edit", 302)
		return
	}

	targetUser, err := reset_password.ChangeUserPassword(id, editPasswordForm.ResetToken, editPasswordForm.Password)
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Update", c).Info("cannot authenticate reset password")
		InternalServerError(w, r)
		return
	}

	go password_mailer.Changed(targetUser.Email)
	logging.SharedInstance().MethodInfo("PasswordsController", "Update", c).Info("success to change password")
	http.Redirect(w, r, "/sign_in", 302)
	return
}
