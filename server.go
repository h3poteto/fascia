package main
import (
	"flag"
	"net/http"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"./controllers"
)


func Routes(m *web.Mux) {
	m.Get("/sign_in", controllers.CallController(&controllers.Sessions{}, "SignIn"))
	m.Post("/sign_in", controllers.CallController(&controllers.Sessions{}, "NewSession"))
	m.Get("/sign_up", controllers.CallController(&controllers.Registrations{}, "SignUp"))
	m.Post("/sign_up", controllers.CallController(&controllers.Registrations{}, "Registration"))
	m.Get("/stylesheets/*", http.FileServer(http.Dir("./public/assets/")))
	m.Get("/", controllers.CallController(&controllers.Root{}, "Index"))
	m.Post("/projects/", controllers.CallController(&controllers.Projects{}, "Create"))
}

func main() {
	flag.Set("bind", ":9090")
	Routes(goji.DefaultMux)
	goji.Serve()
}
