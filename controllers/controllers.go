package controllers

import (
	"os"
	"fmt"
	"net/http"
	"reflect"
	"github.com/zenazn/goji/web"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	userModel "../models/user"
)

type JsonError struct {
	Error string
}

var githubOauthConf = &oauth2.Config{
	ClientID: os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	Scopes: []string{"repo", "write:repo_hook"},
	Endpoint: github.Endpoint,
}

var cookieStore = sessions.NewCookieStore([]byte("session-kesy"))

func CallController(controller interface{}, action string) interface{} {
	method := reflect.ValueOf(controller).MethodByName(action)
	return method.Interface()
}

func LoginRequired(c web.C, w http.ResponseWriter, r *http.Request) (*userModel.UserStruct, bool) {
	session, err := cookieStore.Get(r, "fascia")
	if err != nil {
		fmt.Printf("cookie error\n")
		return nil, false
	}
	id := session.Values["current_user_id"]
	if id == nil {
		fmt.Printf("not logined\n")
		return nil, false
	}
	current_user, err := userModel.CurrentUser(id.(int64))
	if err != nil {
		fmt.Printf("cannot find login user\n")
		return nil, false
	}
	return current_user, true
}
