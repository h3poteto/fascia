package controllers

import (
	"fmt"
	"net/http"
	"reflect"
	"github.com/zenazn/goji/web"
	"github.com/gorilla/sessions"
	userModel "../models/user"
)

type JsonError struct {
	Error string
}

var cookieStore = sessions.NewCookieStore([]byte("session-kesy"))

func CallController(controller interface{}, action string) interface{} {
	method := reflect.ValueOf(controller).MethodByName(action)
	return method.Interface()
}

func LoginRequired(c web.C, w http.ResponseWriter, r *http.Request) (*userModel.UserStruct, bool) {
	session, err := cookieStore.Get(r, "fascia")
	if err != nil {
		return nil, false
	}
	id := session.Values["current_user_id"]
	if id == nil {
		fmt.Printf("not logined\n")
		return nil, false
	}
	current_user, err := userModel.CurrentUser(id.(int))
	if err != nil {
		fmt.Printf("cannot find login user\n")
		return nil, false
	}
	return &current_user, true
}
