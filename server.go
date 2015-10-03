package main
import (
	"flag"
	"net/http"
	"fmt"
	"syscall"
	"strconv"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/flosch/pongo2"
	_ "github.com/flosch/pongo2-addons"
	"./controllers"
)


func Routes(m *web.Mux) {
	// assets
	m.Get("/stylesheets/*", http.FileServer(http.Dir("./public/assets/")))
	m.Get("/javascripts/*", http.FileServer(http.Dir("./public/assets/")))
	m.Get("/images/*", http.FileServer(http.Dir("./frontend/")))
	m.Get("/fonts/*", http.FileServer(http.Dir("./public/assets/")))
	// routing
	m.Get("/", controllers.CallController(&controllers.Root{}, "Index"))
	m.Get("/sign_in", controllers.CallController(&controllers.Sessions{}, "SignIn"))
	m.Post("/sign_in", controllers.CallController(&controllers.Sessions{}, "NewSession"))
	m.Get("/sign_up", controllers.CallController(&controllers.Registrations{}, "SignUp"))
	m.Post("/sign_up", controllers.CallController(&controllers.Registrations{}, "Registration"))
	m.Get("/auth/github", controllers.CallController(&controllers.Oauth{}, "Github"))
	m.Post("/projects/", controllers.CallController(&controllers.Projects{}, "Create"))
	m.Get("/projects/", controllers.CallController(&controllers.Projects{}, "Index"))
	m.Get("/github/repositories", controllers.CallController(&controllers.Github{}, "Repositories"))
}

func main() {
	pongo2.RegisterFilter("suffixAssetsUpdate", SuffixAssetsUpdate)
	flag.Set("bind", ":9090")
	Routes(goji.DefaultMux)
	goji.Serve()
}

func SuffixAssetsUpdate(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
	assetsFile, ok := in.Interface().(string)
	if !ok {
		return nil, &pongo2.Error{
			Sender: "suffixStylesheet",
			ErrorMsg: fmt.Sprintf("Data must be string %T ('%v')", in, in),
		}
	}

	var file syscall.Stat_t
	syscall.Stat("./public/assets" + assetsFile, &file)
	timestamp, _ := file.Mtim.Unix()
	return pongo2.AsValue(assetsFile + "?update=" + strconv.FormatInt(timestamp, 10)), nil
}
