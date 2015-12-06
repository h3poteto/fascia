package controllers_test

import (
	. "../controllers"
	"../models/user"
	"encoding/json"
	"github.com/flosch/pongo2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

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

func LoginFaker(ts *httptest.Server, email string, password string) int64 {
	CheckCSRFToken = func(r *http.Request, token string) bool { return true }
	id, _ := user.Registration(email, password)
	LoginRequired = func(r *http.Request) (*user.UserStruct, error) {
		current_user, _ := user.CurrentUser(id)
		return current_user, nil
	}
	values := url.Values{}
	values.Add("email", email)
	values.Add("password", password)
	http.PostForm(ts.URL+"/sign_in", values)
	return id
}
