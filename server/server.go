package server

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/server/controllers"
	"github.com/h3poteto/fascia/server/filters"
	"github.com/h3poteto/fascia/server/middlewares"

	"context"
	"fmt"
	"net/http"

	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/echo-contrib/pongor"
	"github.com/flosch/pongo2"
	_ "github.com/flosch/pongo2-addons"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Routes defines all routes
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

	login := e.Group("/")
	login.Use(middlewares.Login())
	sessions := &controllers.Sessions{}
	e.GET("/sign_in", sessions.SignIn)
	e.POST("/sign_in", sessions.NewSession)
	login.POST("session", sessions.Update)
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

	github := &controllers.Github{}
	login.GET("github/repositories", github.Repositories)

	projects := &controllers.Projects{}
	login.POST("projects", projects.Create)
	login.GET("projects", projects.Index)

	p := login.Group("projects")
	p.Use(middlewares.Project())
	p.POST("/:project_id", projects.Update)
	p.GET("/:project_id/show", projects.Show)
	p.POST("/:project_id/fetch_github", projects.FetchGithub)
	p.POST("/:project_id/settings", projects.Settings)
	p.POST("/:project_id/webhook", projects.Webhook)
	p.DELETE("/:project_id", projects.Destroy)

	lists := &controllers.Lists{}
	p.GET("/:project_id/lists", lists.Index)
	p.POST("/:project_id/lists", lists.Create)

	l := p.Group("/:project_id/lists")
	l.Use(middlewares.List())
	l.POST("/:list_id", lists.Update)
	l.POST("/:list_id/hide", lists.Hide)
	l.POST("/:list_id/display", lists.Display)

	tasks := &controllers.Tasks{}
	l.POST("/:list_id/tasks", tasks.Create)

	t := l.Group("/:list_id/tasks")
	t.Use(middlewares.Task())
	t.GET("/:task_id", tasks.Show)
	t.POST("/:task_id/move_task", tasks.MoveTask)
	t.POST("/:task_id", tasks.Update)
	t.DELETE("/:task_id", tasks.Delete)

	listOptions := &controllers.ListOptions{}
	login.GET("list_options", listOptions.Index)

	repositories := &controllers.Repositories{}
	e.POST("/repositories/hooks/github", repositories.Hook)

	// errors
	e.GET("/400", controllers.BadRequest)
	e.GET("/404", controllers.NotFound)
	e.GET("/500", controllers.InternalServerError)
}

// Serve start echo server
func Serve() {
	e := echo.New()
	e.Renderer = PongoRenderer()
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

	e.HTTPErrorHandler = ErrorLogging(e)
	e.Use(middlewares.CustomizeLogger())
	e.Use(middlewares.PanicRecover())
	e.Use(middleware.RequestID())
	Routes(e)

	// Start server in gorutine for graceful shutdown
	go func() {
		s := &http.Server{
			Addr:         ":9090",
			ReadTimeout:  20 * time.Minute,
			WriteTimeout: 20 * time.Minute,
		}
		if err := e.StartServer(s); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

// PongoRenderer prepare pongo2, pongo2filter, and pongor
func PongoRenderer() *pongor.Renderer {
	root := os.Getenv("APPROOT")
	pongo2.RegisterFilter("suffixAssetsUpdate", filters.SuffixAssetsUpdate)
	pongorOption := pongor.PongorOption{
		Directory: filepath.Join(root, "server/templates"),
		Reload:    false,
	}
	return pongor.GetRenderer(pongorOption)
}

type fundamental interface {
	StackTrace() errors.StackTrace
}

// ErrorLogging logging error and call default error handler in echo
func ErrorLogging(e *echo.Echo) func(error, echo.Context) {
	return func(err error, c echo.Context) {
		// pkg/errorsにより生成されたエラーについては，各コントローラで適切にハンドリングすること
		// ここでは予定外のエラーが発生した場合にログを飛ばしたい
		// 予定外のエラーなので，errors.fundamentalとecho.HTTPError以外のエラーだけを拾えれば十分なはずである
		_, isFundamental := err.(fundamental)
		_, isHTTPError := err.(*echo.HTTPError)
		if !isFundamental && !isHTTPError {
			logging.SharedInstance().Controller(c).Error(err)
		}
		e.DefaultHTTPErrorHandler(err, c)
	}
}
