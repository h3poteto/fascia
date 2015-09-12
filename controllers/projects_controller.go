package controllers
import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/zenazn/goji/web"
	"github.com/goji/param"
	projectModel "../models/project"
)

type Projects struct {
}

type NewProjectForm struct {
	Title string `param:"title"`
}

func (u *Projects)Index(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, result := LoginRequired(c, w, r)
	encoder := json.NewEncoder(w)
	if !result {
		error := JsonError{Error: "not logined"}
		encoder.Encode(error)
		return
	}
	projects := current_user.Projects()
	encoder.Encode(projects)
}

func (u *Projects)Create(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, result := LoginRequired(c, w, r)
	encoder := json.NewEncoder(w)
	if !result {
		error := JsonError{Error: "not logined"}
		encoder.Encode(error)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Wrong From", 400)
		return
	}
	var newProjectForm NewProjectForm
	err = param.Parse(r.PostForm, &newProjectForm)

	if err != nil {
		http.Error(w, "Wrong parameter", 500)
		return
	}
	fmt.Printf("post new project parameter: %+v\n", newProjectForm)
	project := projectModel.NewProject(0, current_user.Id, newProjectForm.Title)
	project.Save()
	encoder.Encode(*project)
}
