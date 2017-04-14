package controllers_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/flosch/pongo2"
	"github.com/h3poteto/fascia/db/seed"
	. "github.com/h3poteto/fascia/server/controllers"
	"github.com/h3poteto/fascia/server/filters"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/services"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	AfterEach(func() {
		err := seed.TruncateAll()
		if err != nil {
			panic(err)
		}
	})
	RunSpecs(t, "Controllers Suite")
}

var _ = BeforeSuite(func() {
	pongo2.RegisterFilter("suffixAssetsUpdate", filters.SuffixAssetsUpdate)
	pongo2.DefaultSet = pongo2.NewSet("test", pongo2.MustNewLocalFileSystemLoader("../templates"))
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
	user, err := handlers.RegistrationUser(email, password, password)
	if err != nil {
		panic(err)
	}
	LoginRequired = func(r *http.Request) (*services.User, error) {
		return handlers.FindUser(user.UserEntity.UserModel.ID)
	}
	values := url.Values{}
	values.Add("email", email)
	values.Add("password", password)
	http.PostForm(ts.URL+"/sign_in", values)
	return user.UserEntity.UserModel.ID
}
