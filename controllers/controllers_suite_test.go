package controllers_test

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"net/http/httptest"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/flosch/pongo2"
	. "../controllers"
	"../models/user"

	"testing"

	"../filters"
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controllers Suite")
}

var _ = BeforeSuite(func() {
	pongo2.RegisterFilter("suffixAssetsUpdate", filters.SuffixAssetsUpdate)
	pongo2.DefaultSet = pongo2.NewSet("test", pongo2.MustNewLocalFileSystemLoader("../views"))
})

func ParseJson(res *http.Response) (interface{}, int) {
	defer res.Body.Close()
	var resp interface{}
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(contents, &resp)
	return resp, res.StatusCode
}

func ParseResponse(res *http.Response) (string, int) {
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return string(contents), res.StatusCode
}


func LoginFaker(ts *httptest.Server, email string, password string) {
	CheckCSRFToken = func(r *http.Request, token string) bool { return true }
	id, _ := user.Registration(email, password)
	LoginRequired = func(r *http.Request) (*user.UserStruct, bool) {
		current_user, _ := user.CurrentUser(id)
		return current_user, true
	}
	values := url.Values{}
	values.Add("email", email)
	values.Add("password", password)
	http.PostForm(ts.URL + "/sign_in", values)
}
