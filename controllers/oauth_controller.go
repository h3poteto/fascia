package controllers
import (
	"fmt"
	"net/http"
	"github.com/zenazn/goji/web"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	userModel "../models/user"
)

type Oauth struct {
}

func (u *Oauth) Github(c web.C, w http.ResponseWriter, r *http.Request) {
	// 旧セッションの削除
	session, err := cookieStore.Get(r, "fascia")
	session.Options = &sessions.Options{MaxAge: -1}
	session.Save(r, w)

	code := r.URL.Query().Get("code")
	fmt.Printf("github callback param: %+v\n", code)
	token, err := githubOauthConf.Exchange(oauth2.NoContext, code)
	fmt.Printf("token: %v\n", token.AccessToken)
	if err != nil {
		http.Error(w, "Oauth Token Error", 500)
		return
	}
	// userModelにtokenを保存してログイン完了

	current_user, err := userModel.FindOrCreateGithub(token.AccessToken)
	if err != nil {
		fmt.Printf("cannot find user\n")
		http.Redirect(w, r, "/sign_in", 302)
		return
	}
	fmt.Printf("%+v\n", current_user)
	session, err = cookieStore.Get(r, "fascia")
	session.Options = &sessions.Options{Path: "/", MaxAge: 3600}
	session.Values["current_user_id"] = current_user.Id
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "session error", 500)
		return
	}
	http.Redirect(w, r, "/", 302)
	return
}
