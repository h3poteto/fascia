package list

import (
	"../db"
	"../repository"
	"../task"
	"database/sql"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type List interface {
	Save() bool
}

type ListStruct struct {
	Id        int64
	ProjectId int64
	UserId    int64
	Title     sql.NullString
	ListTasks []*task.TaskStruct
	Color     sql.NullString
	database  db.DB
}

func NewList(id int64, projectID int64, userID int64, title string, color string) *ListStruct {
	if projectID == 0 {
		return nil
	}
	nullTitle := sql.NullString{String: title, Valid: true}
	nullColor := sql.NullString{String: color, Valid: true}
	list := &ListStruct{Id: id, ProjectId: projectID, UserId: userID, Title: nullTitle, Color: nullColor}
	list.Initialize()
	return list
}

func FindList(projectID int64, listID int64) *ListStruct {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var id, projectId, userId int64
	var title, color sql.NullString
	rows, _ := table.Query("select id, project_id, user_id, title, color from lists where id = ? AND project_id = ?;", listID, projectID)
	for rows.Next() {
		err := rows.Scan(&id, &projectId, &userId, &title, &color)
		if err != nil {
			panic(err.Error())
		}
	}
	if id != listID {
		fmt.Printf("cannot find list or project did not contain list: %v\n", listID)
		return nil
	} else {
		list := NewList(id, projectId, userId, title.String, color.String)
		return list
	}

}

func (u *ListStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

func (u *ListStruct) Save(repo *repository.RepositoryStruct, OauthToken *sql.NullString) bool {
	table := u.database.Init()
	defer table.Close()
	tx, _ := table.Begin()
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic: %+v\n", err)
			tx.Rollback()
		}
	}()

	result, err := tx.Exec("insert into lists (project_id, user_id, title, color, created_at) values (?, ?, ?, ?, now());", u.ProjectId, u.UserId, u.Title, u.Color)
	if err != nil {
		fmt.Printf("list save error: %+v\n", err)
		tx.Rollback()
		return false
	}

	if OauthToken != nil && OauthToken.Valid && repo != nil {
		token := OauthToken.String
		label := u.CheckLabelPresent(token, repo)
		if label == nil {
			// そもそも既に存在しているなんてことはあまりないのでは
			label = u.CreateGithubLabel(token, repo)
			if label == nil {
				fmt.Printf("github label create failed\n")
				tx.Rollback()
				return false
			}
		} else {
			label = u.UpdateGithubLabel(token, repo)
			if label == nil {
				fmt.Printf("github label update failed\n")
				tx.Rollback()
				return false
			}
		}
	}
	tx.Commit()
	u.Id, _ = result.LastInsertId()
	return true
}

func (u *ListStruct) Update() bool {
	table := u.database.Init()
	defer table.Close()

	_, err := table.Exec("update lists set project_id = ?, user_id = ?, title = ?, color = ? where id = ?;", u.ProjectId, u.UserId, u.Title, u.Color, u.Id)
	if err != nil {
		return false
	}
	return true
}

func (u *ListStruct) Tasks() []*task.TaskStruct {
	table := u.database.Init()
	defer table.Close()

	rows, _ := table.Query("select id, list_id, user_id, title from tasks where list_id = ? order by display_index;", u.Id)
	var slice []*task.TaskStruct
	for rows.Next() {
		var id, listID, userID int64
		var title sql.NullString
		err := rows.Scan(&id, &listID, &userID, &title)
		if err != nil {
			panic(err.Error())
		}
		if listID == u.Id && title.Valid {
			l := task.NewTask(id, listID, userID, title.String)
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

	label := &github.Label{
		Name:  &u.Title.String,
		Color: &u.Color.String,
	}
	githubLabel, response, err := client.Issues.CreateLabel(repo.Owner.String, repo.Name.String, label)
	fmt.Printf("create label for github response: %+v\n", response)
	if err != nil {
		panic(err.Error())
		return nil
	}
	fmt.Printf("github label created: %+v\n", githubLabel)
	return githubLabel
}

func (u *ListStruct) UpdateGithubLabel(token string, repo *repository.RepositoryStruct) *github.Label {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	label := &github.Label{
		Name:  &u.Title.String,
		Color: &u.Color.String,
	}
	githubLabel, response, err := client.Issues.EditLabel(repo.Owner.String, repo.Name.String, u.Title.String, label)
	fmt.Printf("update label for github response: %+v\n", response)
	if err != nil {
		panic(err.Error())
		return nil
	}
	fmt.Printf("github label updated: %+v\n", githubLabel)
	return githubLabel
}
