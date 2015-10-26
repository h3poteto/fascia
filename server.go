package main
import (
	"flag"
	"net/http"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/flosch/pongo2"
	_ "github.com/flosch/pongo2-addons"
	"./controllers"
	"./filters"
)


func Routes(m *web.Mux) {
	// assets
	m.Get("/stylesheets/*", http.FileServer(http.Dir("./public/assets/")))
	m.Get("/javascripts/*", http.FileServer(http.Dir("./public/assets/")))
	m.Get("/images/*", http.FileServer(http.Dir("./frontend/")))
	m.Get("/fonts/*", http.FileServer(http.Dir("./public/assets/")))
	// routing
	m.Get("/", controllers.CallController(&controllers.Root{}, "Index"))
	m.Get("/projects/:project_id", controllers.CallController(&controllers.Root{}, "Index"))
	m.Get("/sign_in", controllers.CallController(&controllers.Sessions{}, "SignIn"))
	m.Post("/sign_in", controllers.CallController(&controllers.Sessions{}, "NewSession"))
	m.Post("/sign_out", controllers.CallController(&controllers.Sessions{}, "SignOut"))
	m.Get("/sign_up", controllers.CallController(&controllers.Registrations{}, "SignUp"))
	m.Post("/sign_up", controllers.CallController(&controllers.Registrations{}, "Registration"))
	m.Get("/auth/github", controllers.CallController(&controllers.Oauth{}, "Github"))
	m.Post("/projects", controllers.CallController(&controllers.Projects{}, "Create"))
	m.Get("/projects", controllers.CallController(&controllers.Projects{}, "Index"))
	m.Get("/projects/:project_id/show", controllers.CallController(&controllers.Projects{}, "Show"))
	m.Get("/github/repositories", controllers.CallController(&controllers.Github{}, "Repositories"))
	m.Get("/projects/:project_id/lists", controllers.CallController(&controllers.Lists{}, "Index"))
	m.Post("/projects/:project_id/lists", controllers.CallController(&controllers.Lists{}, "Create"))
	m.Get("/projects/:project_id/lists/:list_id/tasks", controllers.CallController(&controllers.Tasks{}, "Index"))
	m.Post("/projects/:project_id/lists/:list_id/tasks", controllers.CallController(&controllers.Tasks{}, "Create"))

}

func main() {
	pongo2.DefaultSet = pongo2.NewSet("default", pongo2.MustNewLocalFileSystemLoader("./views"))
	pongo2.RegisterFilter("suffixAssetsUpdate", filters.SuffixAssetsUpdate)
	flag.Set("bind", ":9090")
	Routes(goji.DefaultMux)
	goji.Serve()
}
