package controllers

import (
	userModel "../models/user"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
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

var cookieStore = sessions.NewCookieStore([]byte("session-kesy"))

// ここテストでstubするために関数ポインタをグローバル変数に代入しておきます．もしインスタンスメソッドではない関数をstubする方法があれば，書き換えて構わない．
var CheckCSRFToken = checkCSRF
var LoginRequired = checkLogin

func CallController(controller interface{}, action string) interface{} {
	method := reflect.ValueOf(controller).MethodByName(action)
	return method.Interface()
}

func checkLogin(r *http.Request) (*userModel.UserStruct, error) {
	session, err := cookieStore.Get(r, Key)
	if err != nil {
		return nil, errors.New("cookie error")
	}
	id := session.Values["current_user_id"]
	if id == nil {
		return nil, errors.New("not logined")
	}
	current_user, err := userModel.CurrentUser(id.(int64))
	if err != nil {
		return nil, errors.New("cannot find login user")
	}
	return current_user, nil
}

func GenerateCSRFToken(c web.C, w http.ResponseWriter, r *http.Request) (string, error) {
	session, err := cookieStore.Get(r, Key)
	if err != nil {
		return "", err
	}

	// 現在時間とソルトからトークンを生成
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	io.WriteString(h, "secret_key_salt")
	token := fmt.Sprintf("%x", h.Sum(nil))
	session.Values["token"] = token

	err = cookieStore.Save(r, w, session)
	if err != nil {
		return "", err
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
