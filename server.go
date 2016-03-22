package main

import (
	"./controllers"
	"./filters"
	"./modules/logging"
	"flag"
	"github.com/flosch/pongo2"
	_ "github.com/flosch/pongo2-addons"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func Routes(m *web.Mux) {
	root := os.Getenv("GOJIROOT")
	// assets
	m.Get("/stylesheets/*", http.FileServer(http.Dir(filepath.Join(root, "public/assets/"))))
	m.Get("/javascripts/*", http.FileServer(http.Dir(filepath.Join(root, "public/assets/"))))
	m.Get("/images/*", http.FileServer(http.Dir(filepath.Join(root, "frontend/"))))
	m.Get("/fonts/*", http.FileServer(http.Dir(filepath.Join(root, "public/assets/"))))
	// routing
	m.Get("/about", controllers.CallController(&controllers.Root{}, "About"))
	m.Get("/", controllers.CallController(&controllers.Root{}, "Index"))
	m.Get("/projects/:project_id", controllers.CallController(&controllers.Root{}, "Index"))
	m.Get("/sign_in", controllers.CallController(&controllers.Sessions{}, "SignIn"))
	m.Post("/sign_in", controllers.CallController(&controllers.Sessions{}, "NewSession"))
	m.Post("/sign_out", controllers.CallController(&controllers.Sessions{}, "SignOut"))
	m.Get("/sign_up", controllers.CallController(&controllers.Registrations{}, "SignUp"))
	m.Post("/sign_up", controllers.CallController(&controllers.Registrations{}, "Registration"))
	m.Get("/auth/github", controllers.CallController(&controllers.Oauth{}, "Github"))
	m.Get("/passwords/new", controllers.CallController(&controllers.Passwords{}, "New"))
	m.Post("/passwords/create", controllers.CallController(&controllers.Passwords{}, "Create"))
	m.Get("/passwords/:id/edit", controllers.CallController(&controllers.Passwords{}, "Edit"))
	m.Post("/passwords/:id/update", controllers.CallController(&controllers.Passwords{}, "Update"))

	m.Post("/projects", controllers.CallController(&controllers.Projects{}, "Create"))
	m.Get("/projects", controllers.CallController(&controllers.Projects{}, "Index"))
	m.Post("/projects/:project_id", controllers.CallController(&controllers.Projects{}, "Update"))
	m.Get("/projects/:project_id/show", controllers.CallController(&controllers.Projects{}, "Show"))
	m.Post("/projects/:project_id/fetch_github", controllers.CallController(&controllers.Projects{}, "FetchGithub"))
	m.Post("/projects/:project_id/settings", controllers.CallController(&controllers.Projects{}, "Settings"))
	m.Post("/projects/:project_id/webhook", controllers.CallController(&controllers.Projects{}, "Webhook"))

	m.Get("/github/repositories", controllers.CallController(&controllers.Github{}, "Repositories"))

	m.Get("/projects/:project_id/lists", controllers.CallController(&controllers.Lists{}, "Index"))
	m.Post("/projects/:project_id/lists", controllers.CallController(&controllers.Lists{}, "Create"))
	m.Post("/projects/:project_id/lists/:list_id", controllers.CallController(&controllers.Lists{}, "Update"))

	m.Get("/projects/:project_id/lists/:list_id/tasks", controllers.CallController(&controllers.Tasks{}, "Index"))
	m.Post("/projects/:project_id/lists/:list_id/tasks", controllers.CallController(&controllers.Tasks{}, "Create"))
	m.Post("/projects/:project_id/lists/:list_id/tasks/:task_id/move_task", controllers.CallController(&controllers.Tasks{}, "MoveTask"))

	m.Get("/list_options", controllers.CallController(&controllers.ListOptions{}, "Index"))

	m.Post("/repositories/hooks/github", controllers.CallController(&controllers.Repositories{}, "Hook"))

	// errors
	m.Get("/400", controllers.BadRequest)
	m.Get("/404", controllers.NotFound)
	m.Get("/500", controllers.InternalServerError)
}

func main() {
	root := os.Getenv("GOJIROOT")
	pongo2.DefaultSet = pongo2.NewSet("default", pongo2.MustNewLocalFileSystemLoader(filepath.Join(root, "views")))
	pongo2.RegisterFilter("suffixAssetsUpdate", filters.SuffixAssetsUpdate)
	flag.Set("bind", ":9090")
	mux := goji.DefaultMux
	mux.Use(PanicRecover)
	Routes(mux)

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

func PanicRecover(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(*c)

		defer func() {
			if err := recover(); err != nil {
				logging.SharedInstance().PanicRecover(requestID).Error(err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
