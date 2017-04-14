package services

import (
	"database/sql"
	"fmt"

	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/server/entities/project"
	"github.com/h3poteto/fascia/server/entities/repository"
	"github.com/h3poteto/fascia/server/entities/task"
	"github.com/h3poteto/fascia/server/models/db"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

// Project has a project entity
type Project struct {
	ProjectEntity *project.Project
	database      *sql.DB
}

// NewProject returns a project service
func NewProject(entity *project.Project) *Project {
	return &Project{
		ProjectEntity: entity,
		database:      db.SharedInstance().Connection,
	}
}

// FindProject search project according to project id
func FindProject(projectID int64) (*Project, error) {
	projectEntity, err := project.Find(projectID)
	if err != nil {
		return nil, err
	}
	return NewProject(projectEntity), nil
}

// FindProjectByRepositoryID search project according to repository id
func FindProjectByRepositoryID(repositoryID int64) ([]*Project, error) {
	projectEntities, err := project.FindByRepositoryID(repositoryID)
	if err != nil {
		return nil, err
	}
	var slice []*Project
	for _, e := range projectEntities {
		p := NewProject(e)
		slice = append(slice, p)
	}
	return slice, nil
}

// CheckOwner check project owner as user
func (p *Project) CheckOwner(userID int64) bool {
	if p.ProjectEntity.ProjectModel.UserID != userID {
		return false
	}
	return true
}

// Create create project and repository, and create webhook and related lists.
func (p *Project) Create(userID int64, title string, description string, repositoryID int, oauthToken sql.NullString) (*project.Project, error) {
	var repoID sql.NullInt64
	if repositoryID != 0 && oauthToken.Valid {
		// repository:projectは1:多なので，repositoryがすでに存在している可能性はある
		// そのため先にselectをかけて存在しない場合のみcreateする
		repo, err := repository.FindByGithubRepoID(int64(repositoryID))
		if err != nil {
			repo, err = repository.CreateRepository(repositoryID, oauthToken.String)
			if err != nil {
				return nil, err
			}
		}
		repoID = sql.NullInt64{Int64: repo.RepositoryModel.ID, Valid: true}
	}

	tx, err := p.database.Begin()
	if err != nil {
		return nil, err
	}
	entity := project.New(0, userID, title, description, repoID, true, true)
	if err := entity.Save(tx); err != nil {
		return nil, err
	}

	// 初期listsの保存は今のところ必須要件なのでtransaction内で行いたい
	// ただしgithubへの同期はtransactionが終わってからで良いので，ここでは行わない
	err = entity.CreateInitialLists(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	p.ProjectEntity = entity
	return entity, nil
}

// Update call update in entity
func (p *Project) Update(title string, description string, showIssues bool, showPullRequests bool) error {
	return p.ProjectEntity.Update(title, description, showIssues, showPullRequests)
}

// CreateWebhook call CreateWebhook if project has repository
func (p *Project) CreateWebhook() error {
	oauthToken, err := p.ProjectEntity.OauthToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s://%s/repositories/hooks/github", config.Element("protocol").(string), config.Element("fqdn"))
	repo, find, err := p.ProjectEntity.Repository()
	if err != nil {
		return err
	}
	if !find {
		return nil
	}
	// すでに存在する場合はupdateを叩く
	hook, err := repo.SearchWebhook(oauthToken, url)
	if err != nil {
		return err
	}
	if hook != nil {
		return repo.UpdateWebhook(oauthToken, url, hook)
	}
	return repo.CreateWebhook(oauthToken, url)
}

// DeleteWebhook call DeleteWebhook if project has repository
func (p *Project) DeleteWebhook() error {
	oauthToken, err := p.ProjectEntity.OauthToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s://%s/repositories/hooks/github", config.Element("protocol").(string), config.Element("fqdn"))
	repo, find, err := p.ProjectEntity.Repository()
	if err != nil {
		return err
	}
	if !find {
		return nil
	}
	hook, err := repo.SearchWebhook(oauthToken, url)
	if err != nil {
		return err
	}
	if hook != nil {
		return repo.DeleteWebhook(oauthToken, hook)
	}
	return nil
}

// FetchGithub fetch all lists and all tasks
func (p *Project) FetchGithub() (bool, error) {
	repo, find, err := p.ProjectEntity.Repository()
	// user自体はgithub連携していても，projectが連携していない可能性もあるのでチェック
	if err != nil {
		return false, err
	}
	if !find {
		return false, nil
	}

	oauthToken, err := p.ProjectEntity.OauthToken()
	if err != nil {
		return false, err
	}

	// listを同期
	labels, err := repo.ListLabels(oauthToken)
	if err != nil {
		return false, err
	}
	err = p.ProjectEntity.ListLoadFromGithub(labels)
	if err != nil {
		return false, err
	}

	openIssues, closedIssues, err := repo.GetGithubIssues(oauthToken)
	if err != nil {
		return false, err
	}

	// taskをすべて同期
	err = p.ProjectEntity.TaskLoadFromGithub(append(openIssues, closedIssues...))
	if err != nil {
		return false, err
	}

	// github側へ同期
	// github側に存在するissueについては，#299でDB内に新規作成されるため，ここではgithub側に存在せず，DB内にのみ存在するタスクをissue化すれば良い
	// ここではprojectとlist両方考慮する必要がある
	rows, err := p.database.Query("select tasks.title, tasks.description, lists.title, lists.color from tasks left join lists on lists.id = tasks.list_id where tasks.project_id = ? and tasks.user_id = ? and tasks.issue_number IS NULL;", p.ProjectEntity.ProjectModel.ID, p.ProjectEntity.ProjectModel.UserID)
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
		label, err := repo.CheckLabelPresent(oauthToken, listTitle.String)
		if err != nil {
			return false, err
		}
		if label == nil {
			label, err = repo.CreateGithubLabel(oauthToken, listTitle.String, listColor.String)
			if err != nil {
				return false, err
			}
		}

		_, err = repo.CreateGithubIssue(oauthToken, title, description, []string{*label.Name})
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

// ApplyIssueChanges apply issue changes to task
func (p *Project) ApplyIssueChanges(body github.IssuesEvent) error {
	// taskが見つからない場合は新規作成するのでエラーハンドリング不要
	fmt.Println("debug log: ", p.ProjectEntity)
	fmt.Println("debug log: ", *body.Issue)
	targetTask, _ := task.FindByIssueNumber(p.ProjectEntity.ProjectModel.ID, *body.Issue.Number)

	// create時点ではlabelsが空の状態でhookが飛んできている場合がある
	// editedの場合であってもwebhookにはchangeだけしか載っておらず，最新の状態は載っていない場合がある
	// そのため一度issueの情報を取得し直す必要がある
	issue, err := p.ProjectEntity.ReacquireIssue(body.Issue)
	if err != nil {
		return err
	}
	switch *body.Action {
	case "opened", "reopened":
		if targetTask == nil {
			err = p.ProjectEntity.CreateNewTask(issue)
		} else {
			err = p.ProjectEntity.ReopenTask(targetTask, issue)
		}
	case "closed", "labeled", "unlabeled", "edited":
		err = p.ProjectEntity.TaskApplyLabel(targetTask, issue)
	}
	return err
}

// ApplyPullRequestChanges apply issue changes to task
func (p *Project) ApplyPullRequestChanges(body github.PullRequestEvent) error {
	// taskが見つからない場合は新規作成するのでエラーハンドリング不要
	targetTask, _ := task.FindByIssueNumber(p.ProjectEntity.ProjectModel.ID, *body.Number)

	// note: もしgithubへのアクセスが増大するようであれば，PullRequestオブジェクトからラベルの付替えを行うように改修する

	oauthToken, err := p.ProjectEntity.OauthToken()
	if err != nil {
		return err
	}

	repo, find, err := p.ProjectEntity.Repository()
	if err != nil {
		return err
	}
	if !find {
		return nil
	}
	// CreateNewTaskをするためには，*github.Issueである必要があるが，ここで取得できるのは*github.PullRequestのみなので，なんとかして変換する必要がある
	// そのため問答無用で一度issueを取得し直す
	issue, err := repo.GetGithubIssue(oauthToken, *body.Number)
	if err != nil {
		return err
	}

	switch *body.Action {
	case "opened", "reopened":
		if targetTask == nil {
			err = p.ProjectEntity.CreateNewTask(issue)
		} else {
			err = p.ProjectEntity.ReopenTask(targetTask, issue)
		}
	case "closed", "labeled", "unlabeled", "edited":
		err = p.ProjectEntity.TaskApplyLabel(targetTask, issue)
	}
	return err
}

// FetchCreatedInitialList fetch initial list to github
func (p *Project) FetchCreatedInitialList() error {
	repo, find, err := p.ProjectEntity.Repository()
	if err != nil {
		return err
	}
	if !find {
		return nil
	}

	oauthToken, err := p.ProjectEntity.OauthToken()
	if err != nil {
		return err
	}

	lists, err := p.ProjectEntity.Lists()
	if err != nil {
		return err
	}
	for _, l := range lists {
		label, err := repo.CheckLabelPresent(oauthToken, l.ListModel.Title.String)
		if err != nil {
			return err
		}
		if label != nil {
			continue
		}
		_, err = repo.CreateGithubLabel(oauthToken, l.ListModel.Title.String, l.ListModel.Color.String)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteLists delete all lists related a project
func (p *Project) DeleteLists() error {
	err := p.ProjectEntity.DeleteLists()
	if err != nil {
		return err
	}
	return nil
}

// Delete delete a project record
func (p *Project) Delete() error {
	err := p.ProjectEntity.Delete()
	if err != nil {
		return err
	}
	return nil
}
