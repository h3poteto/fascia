package project

import (
	"../../config"
	"../../modules/hub"
	"../../modules/logging"
	"../db"
	"../list"
	"../list_option"
	"../repository"
	"database/sql"
	"errors"
	"fmt"
)

type Project interface {
	Lists() []*list.ListStruct
	Save() bool
}

type ProjectStruct struct {
	ID               int64
	UserID           int64
	Title            string
	Description      string
	RepositoryID     sql.NullInt64
	ShowIssues       bool
	ShowPullRequests bool
	database         db.DB
}

func NewProject(id int64, userID int64, title string, description string, repositoryID sql.NullInt64, showIssues bool, showPullRequests bool) *ProjectStruct {
	if userID == 0 {
		return nil
	}

	project := &ProjectStruct{ID: id, UserID: userID, Title: title, Description: description, RepositoryID: repositoryID, ShowIssues: showIssues, ShowPullRequests: showPullRequests}
	project.Initialize()
	return project
}

func FindProject(projectID int64) *ProjectStruct {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var id, userID int64
	var repositoryID sql.NullInt64
	var title string
	var description string
	var showIssues, showPullRequests bool
	err := table.QueryRow("select id, user_id, repository_id, title, description, show_issues, show_pull_requests from projects where id = ?;", projectID).Scan(&id, &userID, &repositoryID, &title, &description, &showIssues, &showPullRequests)
	if err != nil {
		logging.SharedInstance().MethodInfo("Project", "FindProject", true).Errorf("cannot find project: %v", err)
		return nil
	}
	project := NewProject(id, userID, title, description, repositoryID, showIssues, showPullRequests)
	return project
}

func Create(userID int64, title string, description string, repositoryID int64, repositoryOwner string, repositoryName string, oauthToken sql.NullString) (p *ProjectStruct, e error) {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()
	tx, _ := table.Begin()
	defer func() {
		if err := recover(); err != nil {
			logging.SharedInstance().MethodInfo("Project", "Create", true).Error("unexpected error")
			tx.Rollback()
			e = errors.New("unexpected error")
			p = nil
		}
	}()

	var repoID sql.NullInt64
	var repo *repository.RepositoryStruct
	if repositoryID != 0 {
		key := repository.GenerateWebhookKey(repositoryName)
		repo = repository.NewRepository(0, repositoryID, repositoryOwner, repositoryName, key)
		if !repo.Save() {
			tx.Rollback()
			logging.SharedInstance().MethodInfo("Project", "Create", true).Error("failed to save repository")
			return nil, errors.New("repository save error")
		}
		repoID = sql.NullInt64{Int64: repo.ID, Valid: true}
	}

	project := NewProject(0, userID, title, description, repoID, true, true)
	if !project.Save() {
		tx.Rollback()
		logging.SharedInstance().MethodInfo("Project", "Create", true).Error("failed to save project")
		return nil, errors.New("failed to save project")
	}

	// github側にwebhooko登録
	err := project.CreateWebhook()
	if err != nil {
		logging.SharedInstance().MethodInfo("Project", "Create").Infof("failed to create webhook: %v", err)
		err = nil
	}

	// 初期リストの準備
	closeListOption := list_option.FindByAction("close")
	if closeListOption == nil {
		tx.Rollback()
		logging.SharedInstance().MethodInfo("Project", "Create", true).Error("cannot find close list option")
		return nil, errors.New("failed to find close list option")
	}
	todo := list.NewList(0, project.ID, userID, config.Element("init_list").(map[interface{}]interface{})["todo"].(string), "f37b1d", sql.NullInt64{})
	inprogress := list.NewList(0, project.ID, userID, config.Element("init_list").(map[interface{}]interface{})["inprogress"].(string), "5eb95e", sql.NullInt64{})
	done := list.NewList(0, project.ID, userID, config.Element("init_list").(map[interface{}]interface{})["done"].(string), "333333", sql.NullInt64{Int64: closeListOption.ID, Valid: true})
	none := list.NewList(0, project.ID, userID, config.Element("init_list").(map[interface{}]interface{})["none"].(string), "ffffff", sql.NullInt64{})
	if err := none.Save(nil, nil); err != nil {
		tx.Rollback()
		logging.SharedInstance().MethodInfo("Project", "Create", true).Errorf("failed to save none list: %v", err)
		return nil, err
	}

	if project.RepositoryID.Valid {
		if err := todo.Save(repo, &oauthToken); err != nil {
			tx.Rollback()
			logging.SharedInstance().MethodInfo("Project", "Create", true).Errorf("failed to save todo list: %v", err)
			return nil, err
		}
		if err := inprogress.Save(repo, &oauthToken); err != nil {
			tx.Rollback()
			logging.SharedInstance().MethodInfo("Project", "Create", true).Errorf("failed to save inprogress list: %v", err)
			return nil, err
		}
		if err := done.Save(repo, &oauthToken); err != nil {
			tx.Rollback()
			logging.SharedInstance().MethodInfo("Project", "Create", true).Errorf("failed to save done list: %v", err)
			return nil, err
		}
	} else {
		if err := todo.Save(nil, nil); err != nil {
			tx.Rollback()
			logging.SharedInstance().MethodInfo("Project", "Create", true).Errorf("failed to save todo list: %v", err)
			return nil, err
		}
		if err := inprogress.Save(nil, nil); err != nil {
			tx.Rollback()
			logging.SharedInstance().MethodInfo("Project", "Create", true).Errorf("failed to save inprogress list: %v", err)
			return nil, err
		}
		if err := done.Save(nil, nil); err != nil {
			tx.Rollback()
			logging.SharedInstance().MethodInfo("Project", "Create", true).Errorf("failed to save done list: %v", err)
			return nil, err
		}
	}
	tx.Commit()
	return project, nil
}

func (u *ProjectStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

func (u *ProjectStruct) Save() bool {
	table := u.database.Init()
	defer table.Close()

	result, err := table.Exec("insert into projects (user_id, repository_id, title, description, show_issues, show_pull_requests, created_at) values (?, ?, ?, ?, ?, ?, now());", u.UserID, u.RepositoryID, u.Title, u.Description, u.ShowIssues, u.ShowPullRequests)
	if err != nil {
		logging.SharedInstance().MethodInfo("Project", "Save", true).Errorf("failed to save project: %v", err)
		return false
	}
	u.ID, _ = result.LastInsertId()
	return true
}

func (u *ProjectStruct) Update(title string, description string, showIssues bool, showPullRequests bool) bool {
	table := u.database.Init()
	defer table.Close()

	u.Title = title
	u.Description = description
	u.ShowIssues = showIssues
	u.ShowPullRequests = showPullRequests
	_, err := table.Exec("update projects set title = ?, description = ?, show_issues = ?, show_pull_requests = ? where id = ?;", u.Title, u.Description, u.ShowIssues, u.ShowPullRequests, u.ID)
	if err != nil {
		logging.SharedInstance().MethodInfo("Project", "Update", true).Errorf("failed to update project: %v", err)
		return false
	}

	return true
}

func (u *ProjectStruct) Lists() []*list.ListStruct {
	table := u.database.Init()
	defer table.Close()

	var slice []*list.ListStruct
	rows, err := table.Query("select id, project_id, user_id, title, color, list_option_id from lists where project_id = ? and title != ?;", u.ID, config.Element("init_list").(map[interface{}]interface{})["none"].(string))
	if err != nil {
		logging.SharedInstance().MethodInfo("Project", "Lists", true).Panic(err)
		return slice
	}
	for rows.Next() {
		var id, projectID, userID int64
		var title, color sql.NullString
		var optionID sql.NullInt64
		err = rows.Scan(&id, &projectID, &userID, &title, &color, &optionID)
		if err != nil {
			logging.SharedInstance().MethodInfo("Project", "Lists", true).Panic(err)
		}
		if projectID == u.ID && title.Valid {
			l := list.NewList(id, projectID, userID, title.String, color.String, optionID)
			slice = append(slice, l)
		}
	}
	return slice
}

func (u *ProjectStruct) NoneList() *list.ListStruct {
	table := u.database.Init()
	defer table.Close()

	var id, projectID, userID int64
	var title, color sql.NullString
	var optionID sql.NullInt64
	err := table.QueryRow("select id, project_id, user_id, title, color, list_option_id from lists where project_id = ? and title = ?;", u.ID, config.Element("init_list").(map[interface{}]interface{})["none"].(string)).Scan(&id, &projectID, &userID, &title, &color, &optionID)
	if err != nil {
		// noneが存在しないということはProjectsController#Createがうまく行ってないので，そっちでエラーハンドリングしてほしい
		logging.SharedInstance().MethodInfo("Project", "NoneList", true).Panic(err)
	}
	if projectID == u.ID && title.Valid {
		return list.NewList(id, projectID, userID, title.String, color.String, optionID)
	}
	return nil
}

func (u *ProjectStruct) Repository() *repository.RepositoryStruct {
	table := u.database.Init()
	defer table.Close()

	var id, repositoryID int64
	var owner, name sql.NullString
	var webhookKey string
	err := table.QueryRow("select repositories.id, repositories.repository_id, repositories.owner, repositories.name, repositories.webhook_key from projects inner join repositories on repositories.id = projects.repository_id where projects.id = ?;", u.ID).Scan(&id, &repositoryID, &owner, &name, &webhookKey)
	if err != nil {
		logging.SharedInstance().MethodInfo("project", "Repository").Infof("cannot find repository: %v", err)
		return nil
	}
	if id == u.RepositoryID.Int64 && owner.Valid {
		r := repository.NewRepository(id, repositoryID, owner.String, name.String, webhookKey)
		return r
	} else {
		logging.SharedInstance().MethodInfo("project", "Repository", true).Error("repository owner discord from project owner")
		return nil
	}
}

func (u *ProjectStruct) FetchGithub() (bool, error) {
	table := u.database.Init()
	defer table.Close()

	repo := u.Repository()
	// user自体はgithub連携していても，projectが連携していない可能性もあるのでチェック
	if repo == nil {
		logging.SharedInstance().MethodInfo("project", "FetchGithub").Info("repository related project is nil")
		return false, errors.New("project did not related to repository")
	}

	oauthToken, err := u.OauthToken()
	if err != nil {
		logging.SharedInstance().MethodInfo("project", "FetchGithub").Infof("oauth token is required: %v", err)
		return false, err
	}

	openIssues, closedIssues, err := hub.GetGithubIssues(oauthToken, repo)
	if err != nil {
		return false, err
	}

	err = u.LoadFromGithub(append(openIssues, closedIssues...))
	if err != nil {
		return false, err
	}

	// github側へ同期
	rows, err := table.Query("select tasks.title, tasks.description, lists.title, lists.color from tasks left join lists on lists.id = tasks.list_id where tasks.user_id = ? and tasks.issue_number IS NULL;", u.UserID)
	if err != nil {
		logging.SharedInstance().MethodInfo("Project", "FetchGithub", true).Panic(err)
		return false, err
	}
	for rows.Next() {
		var title, description string
		var listTitle, listColor sql.NullString
		err := rows.Scan(&title, &description, &listTitle, &listColor)
		if err != nil {
			return false, err
		}
		label, err := hub.CheckLabelPresent(oauthToken, repo, &listTitle.String)
		if err != nil {
			return false, err
		}
		if label == nil {
			label, err = hub.CreateGithubLabel(oauthToken, repo, &listTitle.String, &listColor.String)
			if err != nil {
				return false, errors.New("cannot create github label")
			}
		}
		// ここcreateだけでなくupdateも考慮したほうが良いのではと思ったが，そもそも現状fasciaにはtaskのupdateアクションがないので，updateされることはありえない．そのため，未実装でも問題はない．
		// todo: task#update実装時にはここも実装すること
		_, err = hub.CreateGithubIssue(oauthToken, repo, []string{*label.Name}, &title, &description)
		if err != nil {
			return false, errors.New("cannot create github issue")
		}
	}

	return true, nil
}

// CreateWebhook call hub.CreateWebhook if project has repository
func (u *ProjectStruct) CreateWebhook() error {
	oauthToken, err := u.OauthToken()
	if err != nil {
		logging.SharedInstance().MethodInfo("project", "CreateWebhook").Infof("oauth token is required: %v", err)
		return err
	}

	url := fmt.Sprintf("%s://%s/repositories/hooks/github", config.Element("protocol").(string), config.Element("fqdn"))
	repo := u.Repository()
	if repo == nil {
		return errors.New("cannot find repository")
	}
	err = hub.CreateWebhook(oauthToken, repo, repo.WebhookKey, url)
	return err
}

// OauthToken get oauth token in users
func (u *ProjectStruct) OauthToken() (string, error) {
	table := u.database.Init()
	defer table.Close()

	var oauthToken sql.NullString
	err := table.QueryRow("select users.oauth_token from projects left join users on users.id = projects.user_id where projects.id = ?;", u.ID).Scan(&oauthToken)
	if err != nil {
		logging.SharedInstance().MethodInfo("project", "OauthToken", true).Errorf("oauth_token select error: %v", err)
		return "", err
	}
	if !oauthToken.Valid {
		logging.SharedInstance().MethodInfo("project", "OauthToken").Info("oauth token is not valid")
		return "", errors.New("oauth token isn't exist")
	}

	return oauthToken.String, nil
}
