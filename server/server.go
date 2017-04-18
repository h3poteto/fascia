package server

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/controllers"
	"github.com/h3poteto/fascia/server/filters"

	"os"
	"path/filepath"

	"github.com/echo-contrib/pongor"
	"github.com/flosch/pongo2"
	_ "github.com/flosch/pongo2-addons"
	"github.com/ipfans/echo-session"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
)

func Routes(e *echo.Echo) {
	rootDir := os.Getenv("APPROOT")
	// robots
	e.File("/robots.txt", filepath.Join(rootDir, "public/robots.txt"))

	// assets
	e.Static("/stylesheets", filepath.Join(rootDir, "public/assets/stylesheets"))
	e.Static("/javascripts", filepath.Join(rootDir, "public/assets/javascripts"))
	e.Static("/images", filepath.Join(rootDir, "public/assets/images"))
	e.Static("/fonts", filepath.Join(rootDir, "public/assets/fonts"))
	// routing
	root := &controllers.Root{}
	e.GET("/about", root.About)
	e.GET("/", root.Index)
	e.GET("/projects/:project_id", root.Index)

	sessions := &controllers.Sessions{}
	e.GET("/sign_in", sessions.SignIn)
	e.POST("/sign_in", sessions.NewSession)
	e.POST("/session", sessions.Update)
	e.POST("/sign_out", sessions.SignOut)

	registrations := &controllers.Registrations{}
	e.GET("/sign_up", registrations.SignUp)
	e.POST("/sign_up", registrations.Registration)

	oauth := &controllers.Oauth{}
	e.GET("/auth/github", oauth.Github)

	passwords := &controllers.Passwords{}
	e.GET("/passwords/new", passwords.New)
	e.POST("/passwords/create", passwords.Create)
	e.GET("/passwords/:id/edit", passwords.Edit)
	e.POST("/passwords/:id/update", passwords.Update)

	// webview
	webviews := &controllers.Webviews{}
	e.GET("/webviews/sign_in", webviews.SignIn)
	e.POST("/webviews/sign_in", webviews.NewSession)
	e.GET("/webviews/callback", webviews.Callback)

	projects := &controllers.Projects{}
	e.POST("/projects", projects.Create)
	e.GET("/projects", projects.Index)
	e.POST("/projects/:project_id", projects.Update)
	e.GET("/projects/:project_id/show", projects.Show)
	e.POST("/projects/:project_id/fetch_github", projects.FetchGithub)
	e.POST("/projects/:project_id/settings", projects.Settings)
	e.POST("/projects/:project_id/webhook", projects.Webhook)
	e.DELETE("/projects/:project_id", projects.Destroy)

	github := &controllers.Github{}
	e.GET("/github/repositories", github.Repositories)

	lists := &controllers.Lists{}
	e.GET("/projects/:project_id/lists", lists.Index)
	e.POST("/projects/:project_id/lists", lists.Create)
	e.POST("/projects/:project_id/lists/:list_id", lists.Update)
	e.POST("/projects/:project_id/lists/:list_id/hide", lists.Hide)
	e.POST("/projects/:project_id/lists/:list_id/display", lists.Display)

	tasks := &controllers.Tasks{}
	e.POST("/projects/:project_id/lists/:list_id/tasks", tasks.Create)
	e.GET("/projects/:project_id/lists/:list_id/tasks/:task_id", tasks.Show)
	e.POST("/projects/:project_id/lists/:list_id/tasks/:task_id/move_task", tasks.MoveTask)
	e.POST("/projects/:project_id/lists/:list_id/tasks/:task_id", tasks.Update)
	e.DELETE("/projects/:project_id/lists/:list_id/tasks/:task_id", tasks.Delete)

	listOptions := &controllers.ListOptions{}
	e.GET("/list_options", listOptions.Index)

	repositories := &controllers.Repositories{}
	e.POST("/repositories/hooks/github", repositories.Hook)

	// errors
	e.GET("/400", controllers.BadRequest)
	e.GET("/404", controllers.NotFound)
	e.GET("/500", controllers.InternalServerError)
}

// Serve start echo server
func Serve() {
	root := os.Getenv("APPROOT")
	pongo2.RegisterFilter("suffixAssetsUpdate", filters.SuffixAssetsUpdate)
	pongorOption := pongor.PongorOption{
		Directory: filepath.Join(root, "server/templates"),
		Reload:    false,
	}
	r := pongor.GetRenderer(pongorOption)

	e := echo.New()
	e.Renderer = r
	fqdn := config.Element("fqdn").(interface{})
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://" + fqdn.(string) + ":9090",
			"http://" + fqdn.(string),
			"https://" + fqdn.(string),
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAcceptEncoding,
		},
	}))
	// TODO: customize logger
	sessionKey := os.Getenv("SECRET")
	store := session.NewCookieStore([]byte(sessionKey))
	e.Use(session.Sessions("fascia", store))
	e.Use(middleware.Logger())
	e.Use(PanicRecover())
	e.Use(middleware.RequestID())
	Routes(e)
	e.Logger.Fatal(e.Start(":9090"))
}

func PanicRecover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					var err error
					switch r := r.(type) {
					case error:
						err = r
					default:
						err = errors.Errorf("%v", r)
					}
					logging.SharedInstance().PanicRecover(c).Error(err)
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}

// もしかしたらerrors.Wrapすることでstacktraceが取れて
// それを使ったエラーハンドリングをしたらまとめて取れるかもしれない
// 要検証
