package controllers
import (
	"fmt"
	"net/http"
	"github.com/zenazn/goji/web"
	"github.com/flosch/pongo2"
)

type Root struct {
}
func (u *Root)Index(c web.C, w http.ResponseWriter, r *http.Request) {
	current_user, result := LoginRequired(r)
	fmt.Printf("current_user: %+v\n", current_user)
	if !result {
		fmt.Printf("login required\n")
		http.Redirect(w, r, "/sign_in", 302)
		return
	}
	tpl, err := pongo2.DefaultSet.FromFile("views/home.html.tpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "Fascia"}, w)
}
