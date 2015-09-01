package main
import (
	"flag"
	"net/http"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/flosch/pongo2"
)

func SignIn(c web.C, w http.ResponseWriter, r *http.Request) {
	tpl, err := pongo2.DefaultSet.FromFile("views/sign_in.html.tpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "SignIn"}, w)
}


func Routes(m *web.Mux) {
	m.Get("/sign_in", SignIn)
}

func main() {
	flag.Set("bind", ":9090")
	Routes(goji.DefaultMux)
	goji.Serve()
}
