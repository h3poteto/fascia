package main
import (
	"flag"
	"fmt"
	"net/http"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/goji/param"
	"github.com/flosch/pongo2"
	userModel "./models/user"
)

type SignUpForm struct {
	Email string `param:"email"`
	Password string `param:"password"`
	PasswordConfirm string `param:"password-confirm"`
}


func SignIn(c web.C, w http.ResponseWriter, r *http.Request) {
	tpl, err := pongo2.DefaultSet.FromFile("views/sign_in.html.tpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "SignIn"}, w)
}

func newSession(c web.C, w http.ResponseWriter, r *http.Request) {
}

func SignUp(c web.C, w http.ResponseWriter, r *http.Request) {
	tpl, err := pongo2.DefaultSet.FromFile("views/sign_up.html.tpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "SignUp"}, w)
}

func Registration(c web.C, w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "No good!", 400)
		return
	}

	var signUpForm SignUpForm
	// Parse url.Values (in this case, r.PostForm) and
	// a pointer to our struct so that param can populate it.
	err = param.Parse(r.PostForm, &signUpForm)
	if err != nil {
		http.Error(w, "Real bad.", 500)
		return
	}
	fmt.Printf("%+v\n", signUpForm)
	_ = &userModel.UserStruct{}

	http.Redirect(w, r, "/sign_in", 301)
}


func Routes(m *web.Mux) {
	m.Get("/sign_in", SignIn)
	m.Post("/sign_in", newSession)
	m.Get("/sign_up", SignUp)
	m.Post("/sign_up", Registration)
}

func main() {
	flag.Set("bind", ":9090")
	Routes(goji.DefaultMux)
	goji.Serve()
}
