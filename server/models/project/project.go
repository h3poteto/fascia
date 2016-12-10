package project

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/hub"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/models/db"
	"github.com/h3poteto/fascia/server/models/list"
	"github.com/h3poteto/fascia/server/models/list_option"
	"github.com/h3poteto/fascia/server/models/repository"
	"github.com/h3poteto/fascia/server/services"

	"database/sql"
	"fmt"
	"runtime"

	"github.com/pkg/errors"
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
	database         *sql.DB
}

func NewProject(id int64, userID int64, title string, description string, repositoryID sql.NullInt64, showIssues bool, showPullRequests bool) *ProjectStruct {
	if userID == 0 {
		return nil
	}

	project := &ProjectStruct{ID: id, UserID: userID, Title: title, Description: description, RepositoryID: repositoryID, ShowIssues: showIssues, ShowPullRequests: showPullRequests}
	project.Initialize()
	return project
}

func FindProject(projectID int64) (*ProjectStruct, error) {
	database := db.SharedInstance().Connection

	var id, userID int64
	var repositoryID sql.NullInt64
	var title string
	var description string
	var showIssues, showPullRequests bool
	err := database.QueryRow("select id, user_id, repository_id, title, description, show_issues, show_pull_requests from projects where id = ?;", projectID).Scan(&id, &userID, &repositoryID, &title, &description, &showIssues, &showPullRequests)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	project := NewProject(id, userID, title, description, repositoryID, showIssues, showPullRequests)
	return project, nil
}

// Create create project and repository, and create webhook and related lists.
// TODO: サービス層に移動すべき内容である
func Create(userID int64, title string, description string, repositoryID int, oauthToken sql.NullString) (p *ProjectStruct, e error) {
	database := db.SharedInstance().Connection
	tx, _ := database.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			p = nil
			switch ty := err.(type) {
			case runtime.Error:
				e = errors.Wrap(ty, "runtime error")
			case string:
				e = errors.New(err.(string))
			default:
				e = errors.New("unexpected error")
			}
		}
	}()

	var repoID sql.NullInt64
	var repo *repository.RepositoryStruct
	if repositoryID != 0 && oauthToken.Valid {
		// TODO: 本来サービスにいるべきコードなので，ここで層が違うserviceを呼び出すことを許可する
		repo, err := services.CreateRepository(repositoryID, oauthToken.String)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		repoID = sql.NullInt64{Int64: repo.ID, Valid: true}
	}

	project := NewProject(0, userID, title, description, repoID, true, true)
	if err := project.save(); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 初期リストの準備
	closeListOption, err := list_option.FindByAction("close")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	todo := list.NewList(0, project.ID, userID, config.Element("init_list").(map[interface{}]interface{})["todo"].(string), "f37b1d", sql.NullInt64{}, false)
	inprogress := list.NewList(0, project.ID, userID, config.Element("init_list").(map[interface{}]interface{})["inprogress"].(string), "5eb95e", sql.NullInt64{}, false)
	done := list.NewList(0, project.ID, userID, config.Element("init_list").(map[interface{}]interface{})["done"].(string), "333333", sql.NullInt64{Int64: closeListOption.ID, Valid: true}, false)
	none := list.NewList(0, project.ID, userID, config.Element("init_list").(map[interface{}]interface{})["none"].(string), "ffffff", sql.NullInt64{}, false)
	if err := none.Save(nil, nil); err != nil {
		tx.Rollback()
		return nil, err
	}

	if project.RepositoryID.Valid {
		if err := todo.Save(repo, &oauthToken); err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := inprogress.Save(repo, &oauthToken); err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := done.Save(repo, &oauthToken); err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		if err := todo.Save(nil, nil); err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := inprogress.Save(nil, nil); err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := done.Save(nil, nil); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()

	// callbacks
	go func(project *ProjectStruct) {
		// create webhook
		err := project.CreateWebhook()
		if err != nil {
			logging.SharedInstance().MethodInfo("Project", "Create").Infof("failed to create webhook: %v", err)
		}
		logging.SharedInstance().MethodInfo("Project", "Create").Info("success to create webhook")
		// Sync github
		_, err = project.Repository()
		if err == nil {
			_, err := project.FetchGithub()
			if err != nil {
				logging.SharedInstance().MethodInfoWithStacktrace("Project", "Create", err).Error(err)
			}
		}
	}(project)

	return project, nil
}

func (u *ProjectStruct) Initialize() {
	u.database = db.SharedInstance().Connection
}

func (u *ProjectStruct) save() error {
	result, err := u.database.Exec("insert into projects (user_id, repository_id, title, description, show_issues, show_pull_requests, created_at) values (?, ?, ?, ?, ?, ?, now());", u.UserID, u.RepositoryID, u.Title, u.Description, u.ShowIssues, u.ShowPullRequests)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	u.ID, _ = result.LastInsertId()
	return nil
}

func (u *ProjectStruct) Update(title string, description string, showIssues bool, showPullRequests bool) error {
	u.Title = title
	u.Description = description
	u.ShowIssues = showIssues
	u.ShowPullRequests = showPullRequests
	_, err := u.database.Exec("update projects set title = ?, description = ?, show_issues = ?, show_pull_requests = ? where id = ?;", u.Title, u.Description, u.ShowIssues, u.ShowPullRequests, u.ID)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}

	return nil
}

// Lists list up lists related a project
func (u *ProjectStruct) Lists() ([]*list.ListStruct, error) {
	var slice []*list.ListStruct
	rows, err := u.database.Query("select id, project_id, user_id, title, color, list_option_id, is_hidden from lists where project_id = ? and title != ?;", u.ID, config.Element("init_list").(map[interface{}]interface{})["none"].(string))
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	for rows.Next() {
		var id, projectID, userID int64
		var title, color sql.NullString
		var optionID sql.NullInt64
		var isHidden bool
		err = rows.Scan(&id, &projectID, &userID, &title, &color, &optionID, &isHidden)
		if err != nil {
			return nil, errors.Wrap(err, "sql select error")
		}
		if projectID == u.ID && title.Valid {
			l := list.NewList(id, projectID, userID, title.String, color.String, optionID, isHidden)
			slice = append(slice, l)
		}
	}
	return slice, nil
}

func (u *ProjectStruct) NoneList() (*list.ListStruct, error) {
	var id, projectID, userID int64
	var title, color sql.NullString
	var optionID sql.NullInt64
	var isHidden bool
	err := u.database.QueryRow("select id, project_id, user_id, title, color, list_option_id, is_hidden from lists where project_id = ? and title = ?;", u.ID, config.Element("init_list").(map[interface{}]interface{})["none"].(string)).Scan(&id, &projectID, &userID, &title, &color, &optionID, &isHidden)
	if err != nil {
		// noneが存在しないということはProjectsController#Createがうまく行ってないので，そっちでエラーハンドリングしてほしい
		return nil, errors.Wrap(err, "sql select error")
	}
	if projectID == u.ID && title.Valid {
		return list.NewList(id, projectID, userID, title.String, color.String, optionID, isHidden), nil
	}
	return nil, errors.New("none list not found")
}

func (u *ProjectStruct) Repository() (*repository.RepositoryStruct, error) {
	var id, repositoryID int64
	var owner, name sql.NullString
	var webhookKey string
	err := u.database.QueryRow("select repositories.id, repositories.repository_id, repositories.owner, repositories.name, repositories.webhook_key from projects inner join repositories on repositories.id = projects.repository_id where projects.id = ?;", u.ID).Scan(&id, &repositoryID, &owner, &name, &webhookKey)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	if id == u.RepositoryID.Int64 && owner.Valid {
		r := repository.New(id, repositoryID, owner.String, name.String, webhookKey)
		return r, nil
	}
	return nil, errors.New("repository not found")
}

func (u *ProjectStruct) FetchGithub() (bool, error) {
	repo, err := u.Repository()
	// user自体はgithub連携していても，projectが連携していない可能性もあるのでチェック
	if err != nil {
		return false, err
	}

	oauthToken, err := u.OauthToken()
	if err != nil {
		return false, err
	}

	// listを同期
	labels, err := hub.ListLabels(oauthToken, repo)
	if err != nil {
		return false, err
	}
	err = u.ListLoadFromGithub(labels)
	if err != nil {
		return false, err
	}

	openIssues, closedIssues, err := hub.GetGithubIssues(oauthToken, repo)
	if err != nil {
		return false, err
	}

	// taskをすべて同期
	err = u.TaskLoadFromGithub(append(openIssues, closedIssues...))
	if err != nil {
		return false, err
	}

	// github側へ同期
	// github側に存在するissueについては，#299でDB内に新規作成されるため，ここではgithub側に存在せず，DB内にのみ存在するタスクをissue化すれば良い
	// ここではprojectとlist両方考慮する必要がある
	rows, err := u.database.Query("select tasks.title, tasks.description, lists.title, lists.color from tasks left join lists on lists.id = tasks.list_id where tasks.project_id = ? and tasks.user_id = ? and tasks.issue_number IS NULL;", u.ID, u.UserID)
	if err != nil {
		return false, errors.Wrap(err, "sql select error")
	}
	for rows.Next() {
		var title, description string
		var listTitle, listColor sql.NullString
		err := rows.Scan(&title, &description, &listTitle, &listColor)
		if err != nil {
			return false, errors.Wrap(err, "sql scan error")
		}
		label, err := hub.CheckLabelPresent(oauthToken, repo, &listTitle.String)
		if err != nil {
			return false, err
		}
		if label == nil {
			label, err = hub.CreateGithubLabel(oauthToken, repo, &listTitle.String, &listColor.String)
			if err != nil {
				return false, err
			}
		}

		_, err = hub.CreateGithubIssue(oauthToken, repo, []string{*label.Name}, &title, &description)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

// CreateWebhook call hub.CreateWebhook if project has repository
func (u *ProjectStruct) CreateWebhook() error {
	oauthToken, err := u.OauthToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s://%s/repositories/hooks/github", config.Element("protocol").(string), config.Element("fqdn"))
	repo, err := u.Repository()
	if err != nil {
		return err
	}
	return hub.CreateWebhook(oauthToken, repo, repo.WebhookKey, url)
}

// OauthToken get oauth token in users
func (u *ProjectStruct) OauthToken() (string, error) {
	var oauthToken sql.NullString
	err := u.database.QueryRow("select users.oauth_token from projects left join users on users.id = projects.user_id where projects.id = ?;", u.ID).Scan(&oauthToken)
	if err != nil {
		return "", errors.Wrap(err, "sql select error")
	}
	if !oauthToken.Valid {
		return "", errors.New("oauth token isn't exist")
	}

	return oauthToken.String, nil
}
