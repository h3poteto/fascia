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

func (u *Tasks)Index(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, result := LoginRequired(r)
	encoder := json.NewEncoder(w)
	if !result {
		error := JsonError{Error: "not logined"}
		encoder.Encode(error)
		return
	}
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil && parentProject.UserId.Int64 != current_user.Id {
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
	tasks := parentList.Tasks()
	encoder.Encode(tasks)
	return
}

func (u *Tasks)Create(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, result := LoginRequired(r)
	encoder := json.NewEncoder(w)
	if !result {
		error := JsonError{Error: "not logined"}
		encoder.Encode(error)
		return
	}
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil && parentProject.UserId.Int64 != current_user.Id {
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

	// github同期処理
	repo := parentProject.Repository()
	if current_user.OauthToken.Valid && repo != nil {
		token := current_user.OauthToken.String
		label := parentList.CheckLabelPresent(token, repo)
		// もしラベルがなかった場合は作っておく
		// 色が違っていてもアップデートは不要，それは編集でやってくれ
		if label == nil {
			label = parentList.CreateGithubLabel(token, repo)
			if label == nil {
				error := JsonError{Error: "failed create github label"}
				encoder.Encode(error)
				return
			}
		}
		// issueを作る
		task.CreateGithubIssue(token, repo, []string{parentList.Title.String})
	}
	if !task.Save() {
		error := JsonError{Error: "save failed"}
		encoder.Encode(error)
		return
	}
	encoder.Encode(*task)
}
