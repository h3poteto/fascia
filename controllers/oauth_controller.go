package controllers

import (
	userModel "../models/user"
	"../modules/logging"
	"github.com/gorilla/sessions"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
	"net/http"
)

type Oauth struct {
}

func (u *Oauth) Github(c web.C, w http.ResponseWriter, r *http.Request) {
	// 旧セッションの削除
	session, err := cookieStore.Get(r, "fascia")
	session.Options = &sessions.Options{MaxAge: -1}
	session.Save(r, w)

	code := r.URL.Query().Get("code")
	logging.SharedInstance().MethodInfo("OauthController", "Github").Debugf("github callback param: %+v", code)
	token, err := githubOauthConf.Exchange(oauth2.NoContext, code)
	logging.SharedInstance().MethodInfo("OautController", "Github").Debugf("token: %v", token)
	if err != nil {
		logging.SharedInstance().MethodInfo("OauthController", "Github").Errorf("oauth token error: %v", err)
		http.Error(w, "Oauth Token Error", 500)
		return
	}
	// userModelにtokenを保存してログイン完了
	current_user, err := userModel.FindOrCreateGithub(token.AccessToken)
	if err != nil {
		logging.SharedInstance().MethodInfo("OauthController", "Github").Errorf("cannot find user: %v", err)
		http.Redirect(w, r, "/sign_in", 302)
		return
	}
	logging.SharedInstance().MethodInfo("OauthController", "Github").Debugf("login success: %+v", current_user)
	session, err = cookieStore.Get(r, "fascia")
	session.Options = &sessions.Options{Path: "/", MaxAge: 3600}
	session.Values["current_user_id"] = current_user.Id
	err = session.Save(r, w)
	if err != nil {
		logging.SharedInstance().MethodInfo("OauthController", "Github").Errorf("session error: %v", err)
		http.Error(w, "session error", 500)
		return
	}
	http.Redirect(w, r, "/", 302)
	return
}
