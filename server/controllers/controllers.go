package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	userModel "github.com/h3poteto/fascia/server/models/user"

	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const Key = "fascia"

type JsonError struct {
	Error string
}

var githubOauthConf = &oauth2.Config{
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	Scopes:       []string{"repo", "write:repo_hook", "user:email"},
	Endpoint:     github.Endpoint,
}

var cookieStore = sessions.NewCookieStore([]byte("session-keys"))

// ここテストでstubするために関数ポインタをグローバル変数に代入しておきます．もしインスタンスメソッドではない関数をstubする方法があれば，書き換えて構わない．
var CheckCSRFToken = checkCSRF
var LoginRequired = CheckLogin

func CheckLogin(r *http.Request) (*userModel.UserStruct, error) {
	session, err := cookieStore.Get(r, Key)
	if err != nil {
		return nil, errors.New("cookie error")
	}
	id := session.Values["current_user_id"]
	if id == nil {
		return nil, errors.New("not logined")
	}
	currentUser, err := userModel.CurrentUser(id.(int64))
	if err != nil {
		return nil, err
	}
	return currentUser, nil
}

func GenerateCSRFToken(c web.C, w http.ResponseWriter, r *http.Request) (string, error) {
	session, err := cookieStore.Get(r, Key)
	if err != nil {
		return "", errors.Wrap(err, "cookie error")
	}

	// 現在時間とソルトからトークンを生成
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	io.WriteString(h, "secret_key_salt")
	token := fmt.Sprintf("%x", h.Sum(nil))
	session.Values["token"] = token

	err = cookieStore.Save(r, w, session)
	if err != nil {
		return "", errors.Wrap(err, "cookie error")
	}
	return token, nil
}

func checkCSRF(r *http.Request, token string) bool {
	session, err := cookieStore.Get(r, Key)
	if err != nil {
		return false
	}

	if session.Values["token"] != token {
		return false
	}
	return true
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	tpl, err := pongo2.DefaultSet.FromFile("400.html.tpl")
	if err != nil {
		err := errors.Wrap(err, "template error")
		logging.SharedInstance().MethodInfoWithStacktrace("Controllers", "BadRequest", err).Error(err)
		http.Error(w, "400 BadRequest", 400)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "BadRequest"}, w)
	return
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	tpl, err := pongo2.DefaultSet.FromFile("404.html.tpl")
	if err != nil {
		err := errors.Wrap(err, "template error")
		logging.SharedInstance().MethodInfoWithStacktrace("Controllers", "NotFound", err).Error(err)
		http.Error(w, "404 NotFound", 404)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "NotFound"}, w)
	return
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	tpl, err := pongo2.DefaultSet.FromFile("500.html.tpl")
	if err != nil {
		err := errors.Wrap(err, "template error")
		logging.SharedInstance().MethodInfoWithStacktrace("Controllers", "InternalServerError", err).Error(err)
		http.Error(w, "InternalServerError", 500)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "InternalServerError"}, w)
	return
}
