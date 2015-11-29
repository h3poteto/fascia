package project

import (
	"../../modules/hub"
	"../db"
	"../list"
	"../repository"
	"../task"
	"database/sql"
)

type Project interface {
	Lists() []*list.ListStruct
	Save() bool
}

type ProjectStruct struct {
	Id          int64
	UserId      sql.NullInt64
	Title       string
	Description string
	database    db.DB
}

func NewProject(id int64, userID int64, title string, description string) *ProjectStruct {
	if userID == 0 {
		return nil
	}
	nullUserID := sql.NullInt64{Int64: int64(userID), Valid: true}
	project := &ProjectStruct{Id: id, UserId: nullUserID, Title: title, Description: description}
	project.Initialize()
	return project
}

func FindProject(projectID int64) *ProjectStruct {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var id int64
	var userID sql.NullInt64
	var title string
	var description string
	rows, _ := table.Query("select id, user_id, title, description from projects where id = ?;", projectID)
	for rows.Next() {
		err := rows.Scan(&id, &userID, &title, &description)
		if err != nil {
			panic(err.Error())
		}
	}
	if userID.Valid {
		project := NewProject(id, userID.Int64, title, description)
		return project
	} else {
		return nil
	}
}

func (u *ProjectStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

func (u *ProjectStruct) Save() bool {
	table := u.database.Init()
	defer table.Close()

	result, err := table.Exec("insert into projects (user_id, title, description, created_at) values (?, ?, ?, now());", u.UserId, u.Title, u.Description)
	if err != nil {
		return false
	}
	u.Id, _ = result.LastInsertId()
	return true
}

func (u *ProjectStruct) Lists() []*list.ListStruct {
	table := u.database.Init()
	defer table.Close()

	rows, _ := table.Query("select id, project_id, user_id, title, color from lists where project_id = ?;", u.Id)
	var slice []*list.ListStruct
	for rows.Next() {
		var id, projectID, userID int64
		var title, color sql.NullString
		err := rows.Scan(&id, &projectID, &userID, &title, &color)
		if err != nil {
			panic(err.Error())
		}
		if projectID == u.Id && title.Valid {
			l := list.NewList(id, projectID, userID, title.String, color.String)
			slice = append(slice, l)
		}
	}
	return slice
}

func (u *ProjectStruct) Repository() *repository.RepositoryStruct {
	table := u.database.Init()
	defer table.Close()

	rows, _ := table.Query("select id, project_id, repository_id, owner, name from repositories where project_id = ?", u.Id)
	for rows.Next() {
		var id, projectId, repositoryId int64
		var owner, name sql.NullString
		err := rows.Scan(&id, &projectId, &repositoryId, &owner, &name)
		if err != nil {
			panic(err.Error())
		}
		if projectId == u.Id && owner.Valid {
			r := repository.NewRepository(id, projectId, repositoryId, owner.String, name.String)
			return r
		}
	}
	return nil
}

func (u *ProjectStruct) FetchGithub() bool {
	table := u.database.Init()
	defer table.Close()

	var oauthToken sql.NullString
	err := table.QueryRow("select users.oauth_token from projects left join users on users.id = projects.user_id where projects.id = ?;", u.Id).Scan(&oauthToken)
	if err != nil {
		panic(err.Error())
	}
	if !oauthToken.Valid {
		return false
	}
	openIssues, closedIssues := hub.GetGithubIssues(oauthToken.String, u.Repository())

	for _, issue := range append(openIssues, closedIssues...) {
		var githubLabels []list.ListStruct
		var openList, closedList *list.ListStruct
		for _, label := range issue.Labels {
			for _, list := range u.Lists() {
				// openとcloseのリストは用意しておく
				if list.Title.Valid && list.Title.String == "ToDo" {
					openList = list
				} else if list.Title.Valid && list.Title.String == "Done" {
					closedList = list
				}
				// 紐付いているlabelのlistを持っている時
				if list.Title.Valid && list.Title.String == *label.Name {
					githubLabels = append(githubLabels, *list)
				}
			}
		}
		issueTask := task.FindByIssueNumber(*issue.Number)
		if issueTask == nil {
			issueTask = task.NewTask(0, 0, u.UserId.Int64, sql.NullInt64{Int64: int64(*issue.Number), Valid: true}, *issue.Title)
		}
		if len(githubLabels) == 1 {
			// 一つのlistだけが該当するとき
			issueTask.ListId = githubLabels[0].Id
		} else if len(githubLabels) > 1 {
			// 複数のlistが該当するとき
			issueTask.ListId = githubLabels[0].Id
		} else {
			// ついているlabelのlistを持ってない時
			if *issue.State == "open" && openList != nil {
				issueTask.ListId = openList.Id
			} else if closedList != nil {
				issueTask.ListId = closedList.Id
			}
		}
		// ここはgithub側への同期不要
		if issueTask.Id == 0 {
			issueTask.Save(nil, nil)
		} else {
			issueTask.Update(nil, nil)
		}
	}

	return true
}
