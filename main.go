package main
import (
	"flag"
	"net/http"
	"html/template"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func SignIn(c web.C, w http.ResponseWriter, r *http.Request) {
	template := template.Must(template.ParseFiles("views/sign_in.html"))
	body := "hoge"
	template.Execute(w, body)
}

func Routes(m *web.Mux) {
	m.Get("/sign_in", SignIn)
}

func main() {
	flag.Set("bind", ":9090")
	Routes(goji.DefaultMux)
	goji.Serve()
}
