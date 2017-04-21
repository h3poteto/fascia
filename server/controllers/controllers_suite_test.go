package controllers_test

import (
	"testing"

	"github.com/flosch/pongo2"
	"github.com/h3poteto/fascia/db/seed"
	. "github.com/h3poteto/fascia/server/controllers"
	"github.com/h3poteto/fascia/server/filters"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/services"
	"github.com/labstack/echo"
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

func LoginFaker(email string, password string) int64 {
	CheckCSRFToken = func(c echo.Context, token string) bool { return true }
	user, err := handlers.RegistrationUser(email, password, password)
	if err != nil {
		panic(err)
	}
	LoginRequired = func(c echo.Context) (*services.User, error) {
		return handlers.FindUser(user.UserEntity.UserModel.ID)
	}
	return user.UserEntity.UserModel.ID
}
