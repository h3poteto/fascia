package main
import (
	"flag"
	"fmt"
	"net/http"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/goji/param"
	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
	userModel "./models/user"
	"./controllers"
)

var cookieStore = sessions.NewCookieStore([]byte("session-kesy"))

type SignUpForm struct {
	Email string `param:"email"`
	Password string `param:"password"`
	PasswordConfirm string `param:"password-confirm"`
}

type SignInForm struct {
	Email string `param:"email"`
	Password string `param:"password"`
}

func Root(c web.C, w http.ResponseWriter, r *http.Request) {
	current_user, result := LoginRequired(c, w, r)
	if result {
		fmt.Printf("current_user: %+v\n", *current_user)
	}
	tpl, err := pongo2.DefaultSet.FromFile("views/home.html.tpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "Fascia"}, w)
}

func SignIn(c web.C, w http.ResponseWriter, r *http.Request) {
	tpl, err := pongo2.DefaultSet.FromFile("views/sign_in.html.tpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "SignIn"}, w)
}

func LoginRequired(c web.C, w http.ResponseWriter, r *http.Request) (*userModel.UserStruct, bool) {
	session, err := cookieStore.Get(r, "fascia")
	if err != nil {
		return nil, false
	}
	id := session.Values["current_user_id"]
	if id == nil {
		http.Redirect(w, r, "/sign_in", 301)
		return nil, false
	}
	current_user, err := userModel.CurrentUser(id.(int))
	if err != nil {
		return nil, false
	}
	return &current_user, true
}

func newSession(c web.C, w http.ResponseWriter, r *http.Request) {
	// 旧セッションの削除
	session, err := cookieStore.Get(r, "fascia")
	session.Options = &sessions.Options{MaxAge: -1}
	session.Save(r, w)
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "No good!", 400)
		return
	}
	var signInForm SignInForm
	err = param.Parse(r.PostForm, &signInForm)
	if err != nil {
		http.Error(w, "Real bad", 500)
		return
	}

	current_user, err := userModel.Login(signInForm.Email, signInForm.Password)
	if err != nil {
		http.Redirect(w, r, "/sign_in", 301)
	}
	fmt.Printf("%+v\n", current_user)
	session, err = cookieStore.Get(r, "fascia")
	session.Options = &sessions.Options{MaxAge: 3600}
	session.Values["current_user_id"] = current_user.Id
	session.Save(r, w)
	http.Redirect(w, r, "/", 301)
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
	if signUpForm.Password == signUpForm.PasswordConfirm {
		// login
		res := userModel.Registration(signUpForm.Email, signUpForm.Password)
		if !res {
			http.Redirect(w, r, "/sign_up", 301)
		} else {
			http.Redirect(w, r, "/sign_in", 301)
		}
	} else {
		// error
		http.Redirect(w, r, "/sign_up", 301)
	}
}


func Routes(m *web.Mux) {
	m.Get("/sign_in", SignIn)
	m.Post("/sign_in", newSession)
	m.Get("/sign_up", SignUp)
	m.Post("/sign_up", Registration)
	m.Get("/stylesheets/*", http.FileServer(http.Dir("./public/assets/")))
	m.Get("/", controllers.RootController(controllers.Index))
}

func main() {
	flag.Set("bind", ":9090")
	Routes(goji.DefaultMux)
	goji.Serve()
}
