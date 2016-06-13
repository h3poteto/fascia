package controllers

import (
	"../config"
	userModel "../models/user"
	"../modules/logging"

	"net/http"

	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
)

type Oauth struct {
}

func (u *Oauth) Github(c web.C, w http.ResponseWriter, r *http.Request) {
	// 旧セッションの削除
	session, err := cookieStore.Get(r, "fascia")
	session.Options = &sessions.Options{MaxAge: -1}
	session.Save(r, w)

	code := r.URL.Query().Get("code")
	logging.SharedInstance().MethodInfo("OauthController", "Github", c).Debugf("github callback param: %+v", code)
	token, err := githubOauthConf.Exchange(oauth2.NoContext, code)
	logging.SharedInstance().MethodInfo("OautController", "Github", c).Debugf("token: %v", token)
	if err != nil {
		err := errors.Wrap(err, "oauth token error")
		logging.SharedInstance().MethodInfoWithStacktrace("OauthController", "Github", err, c).Error(err)
		http.Error(w, "Oauth Token Error", 500)
		return
	}
	// userModelにtokenを保存してログイン完了
	currentUser, err := userModel.FindOrCreateGithub(token.AccessToken)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("OauthController", "Github", err, c).Error(err)
		http.Redirect(w, r, "/sign_in", 302)
		return
	}
	logging.SharedInstance().MethodInfo("OauthController", "Github", c).Debugf("login success: %+v", currentUser)
	session, err = cookieStore.Get(r, "fascia")
	session.Options = &sessions.Options{Path: "/", MaxAge: config.Element("session").(map[interface{}]interface{})["timeout"].(int)}
	session.Values["current_user_id"] = currentUser.ID
	err = session.Save(r, w)
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().MethodInfoWithStacktrace("OauthController", "Github", err, c).Error(err)
		http.Error(w, "session error", 500)
		return
	}
	logging.SharedInstance().MethodInfo("OauthController", "Github", c).Info("github login success")

	// iosからのセッションの場合はリダイレクト先を変更
	cookie, err := r.Cookie("fascia-ios")
	if err == nil && cookie.Value == "login-session" {
		http.Redirect(w, r, "/webviews/callback", 302)
		return
	}
	http.Redirect(w, r, "/", 302)
	return
}
