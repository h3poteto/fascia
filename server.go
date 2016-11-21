package main

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/controllers"
	"github.com/h3poteto/fascia/filters"
	"github.com/h3poteto/fascia/modules/logging"

	"flag"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/flosch/pongo2"
	_ "github.com/flosch/pongo2-addons"
	"github.com/goji/glogrus"
	"github.com/rs/cors"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

//go:generate go-bindata -ignore=\\.go -o=config/bindata.go -pkg=config -prefix=config/ config/

func Routes(m *web.Mux) {
	rootDir := os.Getenv("GOJIROOT")
	// robots
	m.Get("/robots.txt", http.FileServer(http.Dir(filepath.Join(rootDir, "public/"))))

	// assets
	m.Get("/stylesheets/*", http.FileServer(http.Dir(filepath.Join(rootDir, "public/assets/"))))
	m.Get("/javascripts/*", http.FileServer(http.Dir(filepath.Join(rootDir, "public/assets/"))))
	m.Get("/images/*", http.FileServer(http.Dir(filepath.Join(rootDir, "frontend/"))))
	m.Get("/fonts/*", http.FileServer(http.Dir(filepath.Join(rootDir, "public/assets/"))))
	// routing
	root := &controllers.Root{}
	m.Get("/about", root.About)
	m.Get("/", root.Index)
	m.Get("/projects/:project_id", root.Index)

	sessions := &controllers.Sessions{}
	m.Get("/sign_in", sessions.SignIn)
	m.Post("/sign_in", sessions.NewSession)
	m.Post("/session", sessions.Update)
	m.Post("/sign_out", sessions.SignOut)

	registrations := &controllers.Registrations{}
	m.Get("/sign_up", registrations.SignUp)
	m.Post("/sign_up", registrations.Registration)

	oauth := &controllers.Oauth{}
	m.Get("/auth/github", oauth.Github)

	passwords := &controllers.Passwords{}
	m.Get("/passwords/new", passwords.New)
	m.Post("/passwords/create", passwords.Create)
	m.Get("/passwords/:id/edit", passwords.Edit)
	m.Post("/passwords/:id/update", passwords.Update)

	// webview
	webviews := &controllers.Webviews{}
	m.Get("/webviews/sign_in", webviews.SignIn)
	m.Post("/webviews/sign_in", webviews.NewSession)
	m.Get("/webviews/callback", webviews.Callback)

	projects := &controllers.Projects{}
	m.Post("/projects", projects.Create)
	m.Get("/projects", projects.Index)
	m.Post("/projects/:project_id", projects.Update)
	m.Get("/projects/:project_id/show", projects.Show)
	m.Post("/projects/:project_id/fetch_github", projects.FetchGithub)
	m.Post("/projects/:project_id/settings", projects.Settings)
	m.Post("/projects/:project_id/webhook", projects.Webhook)

	github := &controllers.Github{}
	m.Get("/github/repositories", github.Repositories)

	lists := &controllers.Lists{}
	m.Get("/projects/:project_id/lists", lists.Index)
	m.Post("/projects/:project_id/lists", lists.Create)
	m.Post("/projects/:project_id/lists/:list_id", lists.Update)
	m.Post("/projects/:project_id/lists/:list_id/hide", lists.Hide)
	m.Post("/projects/:project_id/lists/:list_id/display", lists.Display)

	tasks := &controllers.Tasks{}
	m.Post("/projects/:project_id/lists/:list_id/tasks", tasks.Create)
	m.Get("/projects/:project_id/lists/:list_id/tasks/:task_id", tasks.Show)
	m.Post("/projects/:project_id/lists/:list_id/tasks/:task_id/move_task", tasks.MoveTask)
	m.Post("/projects/:project_id/lists/:list_id/tasks/:task_id", tasks.Update)
	m.Delete("/projects/:project_id/lists/:list_id/tasks/:task_id", tasks.Delete)

	listOptions := &controllers.ListOptions{}
	m.Get("/list_options", listOptions.Index)

	repositories := &controllers.Repositories{}
	m.Post("/repositories/hooks/github", repositories.Hook)

	// errors
	m.Get("/400", controllers.BadRequest)
	m.Get("/404", controllers.NotFound)
	m.Get("/500", controllers.InternalServerError)

	// 任意のファイルも一応ホスティングできるようにしておく
	m.Get("/*", http.FileServer(http.Dir(filepath.Join(rootDir, "public/statics/"))))
}

func main() {
	root := os.Getenv("GOJIROOT")
	pongo2.DefaultSet = pongo2.NewSet("default", pongo2.MustNewLocalFileSystemLoader(filepath.Join(root, "views")))
	pongo2.RegisterFilter("suffixAssetsUpdate", filters.SuffixAssetsUpdate)
	flag.Set("bind", ":9090")
	mux := goji.DefaultMux
	mux.Use(PanicRecover)
	Routes(mux)

	fqdn := config.Element("fqdn").(interface{})
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://" + fqdn.(string) + ":9090",
			"http://" + fqdn.(string),
			"https://" + fqdn.(string),
		},
	})
	goji.Use(c.Handler)

	goji.Use(glogrus.NewGlogrus(logging.SharedInstance().Log, "fascia"))
	goji.Abandon(middleware.Logger)

	fd := flag.Uint("fd", 0, "File descriptor to listen and serve.")
	flag.Parse()

	if *fd != 0 {
		listener, err := net.FileListener(os.NewFile(uintptr(*fd), ""))
		if err != nil {
			panic(err)
		}
		goji.ServeListener(listener)
	} else {
		goji.Serve()
	}
}

// PanicRecover recover any panic and send information to logrus
func PanicRecover(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				logging.SharedInstance().PanicRecover(*c).Error(err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
