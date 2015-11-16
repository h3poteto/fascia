package controllers

import (
	listModel "../models/list"
	projectModel "../models/project"
	taskModel "../models/task"
	"encoding/json"
	"fmt"
	"github.com/goji/param"
	"github.com/zenazn/goji/web"
	"net/http"
	"strconv"
)

type Tasks struct {
}

type NewTaskForm struct {
	Title string `param:"title"`
}

type MoveTaskFrom struct {
	ToListId     int64 `param:"to_list_id"`
	prevToTaskId int64 `param:"prev_to_task_id"`
}

func (u *Tasks) Index(c web.C, w http.ResponseWriter, r *http.Request) {
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

func (u *Tasks) Create(c web.C, w http.ResponseWriter, r *http.Request) {
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

func (u *Tasks) MoveTask(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, result := LoginRequired(r)
	encoder := json.NewEncoder(w)
	if !result {
		http.Error(w, "not logined", 401)
		return
	}
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil || parentProject.UserId.Int64 != current_user.Id {
		http.Error(w, "project not found", 404)
		return
	}

	listID, _ := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	parentList := listModel.FindList(parentProject.Id, listID)
	if parentList == nil {
		http.Error(w, "list not found", 404)
		return
	}

	taskID, _ := strconv.ParseInt(c.URLParams["task_id"], 10, 64)
	task := taskModel.FindTask(parentList.Id, taskID)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "WrongForm", 400)
		return
	}
	var moveTaskFrom MoveTaskFrom
	err = param.Parse(r.PostForm, &moveTaskFrom)
	if err != nil {
		http.Error(w, "Wrong parameter", 500)
		return
	}
	fmt.Printf("post move taks parameter: %+v\n", moveTaskFrom)
	if !task.ChangeList(moveTaskFrom.ToListId, moveTaskFrom.prevToTaskId) {
		error := JsonError{Error: "list change failed"}
		encoder.Encode(error)
		return
	}
	allLists := parentProject.Lists()
	for _, l := range allLists {
		l.ListTasks = l.Tasks()
	}
	fmt.Printf("move task: %+v\n", allLists)
	encoder.Encode(allLists)
	return
}
