package controllers_test

import (
	"testing"

	"github.com/flosch/pongo2"
	"github.com/h3poteto/fascia/db/seed"
	"github.com/h3poteto/fascia/server/commands/account"
	"github.com/h3poteto/fascia/server/commands/board"
	. "github.com/h3poteto/fascia/server/controllers"
	"github.com/h3poteto/fascia/server/filters"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/middlewares"
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

func CSRFFaker() {
	CheckCSRFToken = func(c echo.Context, token string) bool { return true }
}

func LoginFaker(c echo.Context, email string, password string) (*account.User, echo.Context) {
	user, err := handlers.LoginUser(email, password)
	if err != nil {
		panic(err)
	}
	var ctx echo.Context
	ctx = &middlewares.LoginContext{
		c,
		user,
	}
	return user, ctx
}

func ProjectContext(c echo.Context, p *board.Project) echo.Context {
	lc, ok := c.(*middlewares.LoginContext)
	if !ok {
		panic("Cast context")
	}
	var ctx echo.Context
	ctx = &middlewares.ProjectContext{
		*lc,
		p,
	}
	return ctx
}

func ListContext(c echo.Context, l *board.List) echo.Context {
	pc, ok := c.(*middlewares.ProjectContext)
	if !ok {
		panic("Cast context")
	}
	var ctx echo.Context
	ctx = &middlewares.ListContext{
		*pc,
		l,
	}
	return ctx
}

func TaskContext(c echo.Context, t *board.Task) echo.Context {
	lc, ok := c.(*middlewares.ListContext)
	if !ok {
		panic("Cast context")
	}
	var ctx echo.Context
	ctx = &middlewares.TaskContext{
		*lc,
		t,
	}
	return ctx
}
