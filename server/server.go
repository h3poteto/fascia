package server

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/server/controllers"
	"github.com/h3poteto/fascia/server/filters"
	"github.com/h3poteto/fascia/server/middlewares"

	"context"
	"net/http"

	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/h3poteto/pongo2echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Routes defines all routes
func Routes(e *echo.Echo) {
	rootDir := os.Getenv("APPROOT")

	jwtConfig := middleware.JWTConfig{
		Skipper:    middlewares.JWTSkipper,
		Claims:     &config.JwtCustomClaims{},
		SigningKey: []byte(os.Getenv("SECRET")),
	}

	// robots
	e.File("/robots.txt", filepath.Join(rootDir, "public/robots.txt"))

	// assets
	e.Static("/lp/css", filepath.Join(rootDir, "public/lp/css"))
	e.Static("/lp/images", filepath.Join(rootDir, "public/lp/images"))
	e.Static("/js", filepath.Join(rootDir, "public/assets/js"))
	e.Static("/fonts", filepath.Join(rootDir, "public/assets/fonts"))
	// routing
	root := &controllers.Root{}
	e.GET("/health_check", root.HealthCheck)
	e.GET("/about", root.About)
	e.GET("/", root.Index)
	e.GET("/projects/:project_id", root.Index)
	e.GET("/projects/:project_id/lists/:list_id/edit", root.Index)
	e.GET("/projects/:project_id/lists/:list_id/tasks/new", root.Index)
	e.GET("/projects/:project_id/lists/:list_id/tasks/:task_id", root.Index)
	e.GET("/projects/:project_id/lists/:list_id/tasks/:task_id/edit", root.Index)
	e.GET("/settings", root.Index)

	sessions := &controllers.Sessions{}
	e.GET("/sign_in", sessions.SignIn)
	e.POST("/sign_in", sessions.Create)
	e.PATCH("/session", sessions.Update, middleware.JWTWithConfig(jwtConfig), middlewares.Login())
	e.GET("/session", sessions.Show, middleware.JWTWithConfig(jwtConfig), middlewares.Login())
	e.DELETE("/sign_out", sessions.SignOut)

	oauth := &controllers.Oauth{}
	e.GET("/oauth/sign_in", oauth.SignIn)
	e.GET("/auth/github", oauth.Github)

	// webview
	webviews := &controllers.Webviews{}
	e.GET("/webviews/oauth/sign_in", webviews.OauthSignIn)
	e.GET("/webviews/callback", webviews.Callback)
	e.GET("/webviews/inquiries/new", webviews.NewInquiry)
	e.POST("/webviews/inquiries", webviews.Inquiry)

	inquiries := &controllers.Inquiries{}
	e.GET("/inquiries/new", inquiries.New)
	e.POST("/inquiries", inquiries.Create)

	github := &controllers.Github{}
	e.GET("/api/github/repositories", github.Repositories, middleware.JWTWithConfig(jwtConfig), middlewares.Login())

	projects := &controllers.Projects{}
	e.POST("/api/projects", projects.Create, middleware.JWTWithConfig(jwtConfig), middlewares.Login())
	e.GET("/api/projects", projects.Index, middleware.JWTWithConfig(jwtConfig), middlewares.Login())

	p := e.Group("/api/projects/:project_id")
	p.Use(middleware.JWTWithConfig(jwtConfig))
	p.Use(middlewares.Login())
	p.Use(middlewares.Project())
	p.PATCH("", projects.Update)
	p.GET("/show", projects.Show)
	p.POST("/fetch_github", projects.FetchGithub)
	p.PATCH("/settings", projects.Settings)
	p.POST("/webhook", projects.Webhook)
	p.DELETE("", projects.Destroy)

	lists := &controllers.Lists{}
	p.GET("/lists", lists.Index)
	p.POST("/lists", lists.Create)

	l := p.Group("/lists/:list_id")
	l.Use(middlewares.List())
	l.GET("", lists.Show)
	l.PATCH("", lists.Update)
	l.PATCH("/hide", lists.Hide)
	l.PATCH("/display", lists.Display)

	tasks := &controllers.Tasks{}
	l.POST("/tasks", tasks.Create)

	t := l.Group("/tasks/:task_id")
	t.Use(middlewares.Task())
	t.GET("", tasks.Show)
	t.POST("/move_task", tasks.MoveTask)
	t.PATCH("", tasks.Update)
	t.DELETE("", tasks.Delete)

	listOptions := &controllers.ListOptions{}
	e.GET("/api/list_options", listOptions.Index, middlewares.Login())

	settings := &controllers.Settings{}
	e.PATCH("/api/settings/password", settings.Password, middlewares.Login())

	repositories := &controllers.Repositories{}
	e.POST("/repositories/hooks/github", repositories.Hook)

	e.GET("/privacy_policy", controllers.PrivacyPolicy)

	e.GET("/*", controllers.NotFound)

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

	e.HTTPErrorHandler = middlewares.ErrorLogging(e)
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

// PongoRenderer prepare pongo2 through pongo2echo
func PongoRenderer() *pongo2echo.Pongo2Echo {
	render := pongo2echo.NewRenderer()
	root := os.Getenv("APPROOT")
	render.RegisterFilter("suffixAssetsUpdate", filters.SuffixAssetsUpdate)
	render.RegisterFilter("digestedAssets", filters.DigestedAssets)
	render.AddDirectory(filepath.Join(root, "server/templates"))

	return render
}
