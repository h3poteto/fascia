package project

import (
	"database/sql"

	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/server/aggregations/list"
	"github.com/h3poteto/fascia/server/aggregations/list_option"
	"github.com/h3poteto/fascia/server/aggregations/repository"
	"github.com/h3poteto/fascia/server/models/db"
	"github.com/h3poteto/fascia/server/models/project"

	"github.com/pkg/errors"
)

type Project struct {
	ProjectModel *project.ProjectStruct
	database     *sql.DB
}

func New(id int64, userID int64, title string, description string, repositoryID sql.NullInt64, showIssues bool, showPullRequests bool) *Project {
	return &Project{
		ProjectModel: project.New(id, userID, title, description, repositoryID, showIssues, showPullRequests),
		database:     db.SharedInstance().Connection,
	}
}

func Find(id int64) (*Project, error) {
	p, err := project.Find(id)
	if err != nil {
		return nil, err
	}
	return &Project{
		ProjectModel: p,
		database:     db.SharedInstance().Connection,
	}, nil
}

func FindByRepositoryID(repositoryID int64) (*Project, error) {
	p, err := project.FindByRepositoryID(repositoryID)
	if err != nil {
		return nil, err
	}
	return &Project{
		ProjectModel: p,
		database:     db.SharedInstance().Connection,
	}, nil
}

func (p *Project) Save(tx *sql.Tx) error {
	return p.ProjectModel.Save(tx)
}

func (p *Project) Update(title string, description string, showIssues bool, showPullRequests bool) error {
	return p.ProjectModel.Update(title, description, showIssues, showPullRequests)
}

func (p *Project) CreateInitialLists(tx *sql.Tx) error {
	// 初期リストの準備
	closeListOption, err := list_option.FindByAction("close")
	if err != nil {
		tx.Rollback()
		return err
	}
	todo := list.New(
		0,
		p.ProjectModel.ID,
		p.ProjectModel.UserID,
		config.Element("init_list").(map[interface{}]interface{})["todo"].(string),
		"f37b1d",
		sql.NullInt64{},
		false,
	)
	inprogress := list.New(
		0,
		p.ProjectModel.ID,
		p.ProjectModel.UserID,
		config.Element("init_list").(map[interface{}]interface{})["inprogress"].(string),
		"5eb95e",
		sql.NullInt64{},
		false,
	)
	done := list.New(
		0,
		p.ProjectModel.ID,
		p.ProjectModel.UserID,
		config.Element("init_list").(map[interface{}]interface{})["done"].(string),
		"333333",
		sql.NullInt64{Int64: closeListOption.ListOptionModel.ID, Valid: true},
		false,
	)
	none := list.New(
		0,
		p.ProjectModel.ID,
		p.ProjectModel.UserID,
		config.Element("init_list").(map[interface{}]interface{})["none"].(string),
		"ffffff",
		sql.NullInt64{},
		false,
	)

	// githubへの同期はもっと上層で行う
	// TODO: ここもtxにしないと駄目
	if err := none.Save(); err != nil {
		tx.Rollback()
		return err
	}

	if err := todo.Save(); err != nil {
		tx.Rollback()
		return err
	}
	if err := inprogress.Save(); err != nil {
		tx.Rollback()
		return err
	}
	if err := done.Save(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// OauthToken get oauth token in users
func (p *Project) OauthToken() (string, error) {
	var oauthToken sql.NullString
	err := p.database.QueryRow("select users.oauth_token from projects left join users on users.id = projects.user_id where projects.id = ?;", p.ProjectModel.ID).Scan(&oauthToken)
	if err != nil {
		return "", errors.Wrap(err, "sql select error")
	}
	if !oauthToken.Valid {
		return "", errors.New("oauth token isn't exist")
	}

	return oauthToken.String, nil
}

// Lists list up lists related a project
func (p *Project) Lists() ([]*list.List, error) {
	var slice []*list.List
	rows, err := p.database.Query("select id, project_id, user_id, title, color, list_option_id, is_hidden from lists where project_id = ? and title != ?;", p.ProjectModel.ID, config.Element("init_list").(map[interface{}]interface{})["none"].(string))
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
		if projectID == p.ProjectModel.ID && title.Valid {
			l := list.New(id, projectID, userID, title.String, color.String, optionID, isHidden)
			slice = append(slice, l)
		}
	}
	return slice, nil
}

func (p *Project) NoneList() (*list.List, error) {
	var id, projectID, userID int64
	var title, color sql.NullString
	var optionID sql.NullInt64
	var isHidden bool
	err := p.database.QueryRow("select id, project_id, user_id, title, color, list_option_id, is_hidden from lists where project_id = ? and title = ?;", p.ProjectModel.ID, config.Element("init_list").(map[interface{}]interface{})["none"].(string)).Scan(&id, &projectID, &userID, &title, &color, &optionID, &isHidden)
	if err != nil {
		// noneが存在しないということはProjectsController#Createがうまく行ってないので，そっちでエラーハンドリングしてほしい
		return nil, errors.Wrap(err, "sql select error")
	}
	if projectID == p.ProjectModel.ID && title.Valid {
		return list.New(id, projectID, userID, title.String, color.String, optionID, isHidden), nil
	}
	return nil, errors.New("none list not found")
}

func (p *Project) Repository() (*repository.Repository, error) {
	var id, repositoryID int64
	var owner, name sql.NullString
	var webhookKey string
	err := p.database.QueryRow("select repositories.id, repositories.repository_id, repositories.owner, repositories.name, repositories.webhook_key from projects inner join repositories on repositories.id = projects.repository_id where projects.id = ?;", p.ProjectModel.ID).Scan(&id, &repositoryID, &owner, &name, &webhookKey)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	if id == p.ProjectModel.RepositoryID.Int64 && owner.Valid {
		r := repository.New(id, repositoryID, owner.String, name.String, webhookKey)
		return r, nil
	}
	return nil, errors.New("repository not found")
}
