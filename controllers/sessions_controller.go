package controllers
import (
	"fmt"
	"net/http"
	"github.com/zenazn/goji/web"
	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
	"github.com/goji/param"
	"golang.org/x/oauth2"
	userModel "../models/user"
)

type Sessions struct {
}

type SignInForm struct {
	Email    string `param:"email"`
	Password string `param:"password"`
	Token    string `param:"token"`
}

func (u *Sessions)SignIn(c web.C, w http.ResponseWriter, r *http.Request) {
	url := githubOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	token, result := GenerateCSRFToken(c, w, r)
	if !result {
		http.Error(w, "Real bad.", 500)
		return
	}

	tpl, err := pongo2.DefaultSet.FromFile("views/sign_in.html.tpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "SignIn", "oauthURL": url, "token": token}, w)
}

func (u *Sessions)NewSession(c web.C, w http.ResponseWriter, r *http.Request) {
	// 旧セッションの削除
	session, err := cookieStore.Get(r, "fascia")
	session.Options = &sessions.Options{MaxAge: -1}
	session.Save(r, w)
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "No good!", 400)
		return
	}
	var signInForm SignInForm
	err = param.Parse(r.PostForm, &signInForm)
	if err != nil {
		http.Error(w, "Real bad", 500)
		return
	}

	if !CheckCSRFToken(r, signInForm.Token) {
		http.Error(w, "Cannot verify CSRF token", 500)
		return
	}

	current_user, err := userModel.Login(signInForm.Email, signInForm.Password)
	if err != nil {
		http.Redirect(w, r, "/sign_in", 301)
		return
	}
	fmt.Printf("%+v\n", current_user)
	session, err = cookieStore.Get(r, "fascia")
	session.Options = &sessions.Options{Path: "/", MaxAge: 3600}
	session.Values["current_user_id"] = current_user.Id
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
	return
}
