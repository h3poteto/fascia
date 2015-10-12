package controllers
import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
	"github.com/zenazn/goji/web"
	"github.com/goji/param"
	projectModel "../models/project"
	listModel "../models/list"
	taskModel "../models/task"
)

type Tasks struct {
}

type NewTaskForm struct {
	Title string `param:"title"`
}

func (u *Tasks)Create(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, result := LoginRequired(c, w, r)
	encoder := json.NewEncoder(w)
	if !result {
		error := JsonError{Error: "not logined"}
		encoder.Encode(error)
		return
	}
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil {
		error := JsonError{Error: "project not found"}
		encoder.Encode(error)
		return
	}
	listID, _ := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	parentList := listModel.FindList(projectID, listID)
	if parentList == nil {
		error := JsonError{Error: "list not found"}
		encoder.Encode(error)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Wrong From", 400)
		return
	}
	var newTaskForm NewTaskForm
	err = param.Parse(r.PostForm, &newTaskForm)
	if err != nil {
		http.Error(w, "Wrong parameter", 500)
		return
	}
	fmt.Printf("post new task parameter: %+v\n", newTaskForm)

	task := taskModel.NewTask(0, parentList.Id, newTaskForm.Title)
	if !task.Save() {
		error := JsonError{Error: "save failed"}
		encoder.Encode(error)
		return
	}
	encoder.Encode(*task)
}
