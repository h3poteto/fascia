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
	current_user, result := LoginRequired(c, w, r)
	if result {
		fmt.Printf("current_user: %+v\n", *current_user)
	} else {
		http.Redirect(w, r, "/sign_in", 301)
		return
	}
	tpl, err := pongo2.DefaultSet.FromFile("views/home.html.tpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "Fascia"}, w)
}
