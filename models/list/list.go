package list

import (
	"fmt"
	"../db"
	"database/sql"
	"../task"
	"../repository"

	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
)

type List interface {
	Save() bool
}

type ListStruct struct {
	Id int64
	ProjectId int64
	Title sql.NullString
	ListTasks []*task.TaskStruct
	Color string
	database db.DB
}

func NewList(id int64, projectID int64, title string) *ListStruct {
	if projectID == 0 {
		return nil
	}
	nullTitle := sql.NullString{String: title, Valid: true}
	list := &ListStruct{Id: id, ProjectId: projectID, Title: nullTitle}
	list.Initialize()
	return list
}

func FindList(projectID int64, listID int64) *ListStruct {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var id, projectId int64
	var title string
	rows, _ := table.Query("select id, project_id, title from lists where id = ? AND project_id = ?;", listID, projectID)
	for rows.Next() {
		err := rows.Scan(&id, &projectId, &title)
		if err != nil {
			panic(err.Error())
		}
	}
	if id != listID {
		fmt.Printf("cannot find list or project did not contain list: %v\n", listID)
		return nil
	} else {
		list := NewList(id, projectId, title)
		return list
	}

}

func (u *ListStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

func (u *ListStruct) Save() bool {
	table := u.database.Init()
	defer table.Close()

	result, err := table.Exec("insert into lists (project_id, title, created_at) values (?, ?, now());", u.ProjectId, u.Title)
	if err != nil {
		return false
	}
	u.Id, _ = result.LastInsertId()
	return true
}

func (u *ListStruct) Update() bool {
	table := u.database.Init()
	defer table.Close()

	_, err := table.Exec("update lists set project_id = ?, title = ? where id = ?;", u.ProjectId, u.Title, u.Id)
	if err != nil {
		return false
	}
	return true
}

func (u *ListStruct) Tasks() []*task.TaskStruct {
	table := u.database.Init()
	defer table.Close()

	rows, _ := table.Query("select id, list_id, title from tasks where list_id = ?;", u.Id)
	var slice []*task.TaskStruct
	for rows.Next() {
		var id, listID int64
		var title sql.NullString
		err := rows.Scan(&id, &listID, &title)
		if err != nil {
			panic(err.Error())
		}
		if listID == u.Id && title.Valid {
			l := task.NewTask(id, listID, title.String)
			slice = append(slice, l)
		}
	}
	return slice
}


func (u *ListStruct) CheckLabelPresent(token string, repo *repository.RepositoryStruct) *github.Label {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	githubLabel, response, err := client.Issues.GetLabel(repo.Owner.String, repo.Name.String, u.Title.String)
	fmt.Printf("get label for github response: %+v\n", response)
	if err != nil {
		fmt.Printf("cannot find github label: %v\n", repo.Name.String)
		return nil
	}
	fmt.Printf("github label: %+v\n", githubLabel)
	return githubLabel
}

func (u *ListStruct) CreateGithubLabel(token string, repo *repository.RepositoryStruct) *github.Label {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	// TODO: listの設定項目に色を追加してほしい．その色を元にここでは色を設定する
	// どうせリストごとに色をつけるのはfascia内の機能としてもほしいので，デフォルトで色はなにか当てておく
	u.Color = "000000"
	label := &github.Label{
		Name: &u.Title.String,
		Color: &u.Color,
	}
	githubLabel, response, err := client.Issues.CreateLabel(repo.Name.String, repo.Owner.String, label)
	fmt.Printf("create label for github response: %+v\n", response)
	if err != nil {
		panic(err.Error())
		return nil
	}
	fmt.Printf("github label created: %+v\n", githubLabel)
	return githubLabel
}
