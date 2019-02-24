package project

import (
	"database/sql"

	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/server/domains/entities/list"
	"github.com/h3poteto/fascia/server/domains/entities/list_option"
	"github.com/h3poteto/fascia/server/infrastructures/project"
)

// Project has a project model object
type Project struct {
	ID               int64
	UserID           int64
	Title            string
	Description      string
	RepositoryID     sql.NullInt64
	ShowIssues       bool
	ShowPullRequests bool
	infrastructure   *project.Project
}

// New returns a project entity
func New(id int64, userID int64, title string, description string, repositoryID sql.NullInt64, showIssues bool, showPullRequests bool) *Project {
	infrastructure := project.New(id, userID, title, description, repositoryID, showIssues, showPullRequests)
	if infrastructure == nil {
		return nil
	}
	p := &Project{
		infrastructure: infrastructure,
	}
	p.reload()
	return p
}

func (p *Project) reflect() {
	p.infrastructure.ID = p.ID
	p.infrastructure.UserID = p.UserID
	p.infrastructure.Title = p.Title
	p.infrastructure.Description = p.Description
	p.infrastructure.RepositoryID = p.RepositoryID
	p.infrastructure.ShowIssues = p.ShowIssues
	p.infrastructure.ShowPullRequests = p.ShowPullRequests
}

func (p *Project) reload() error {
	if p.ID != 0 {
		latestProject, err := project.Find(p.ID)
		if err != nil {
			return err
		}
		p.infrastructure = latestProject
	}
	p.ID = p.infrastructure.ID
	p.UserID = p.infrastructure.UserID
	p.Title = p.infrastructure.Title
	p.Description = p.infrastructure.Description
	p.RepositoryID = p.infrastructure.RepositoryID
	p.ShowIssues = p.infrastructure.ShowIssues
	p.ShowPullRequests = p.infrastructure.ShowPullRequests
	return nil
}

// Save call project model save
func (p *Project) Save(tx *sql.Tx) error {
	p.reflect()
	if err := p.infrastructure.Save(tx); err != nil {
		return err
	}
	return p.reload()
}

// Update call project model update
func (p *Project) Update(title string, description string, showIssues bool, showPullRequests bool) error {
	if err := p.infrastructure.Update(title, description, showIssues, showPullRequests); err != nil {
		return err
	}
	if err := p.reload(); err != nil {
		return err
	}
	return nil
}

// CreateInitialLists create initial lists in self project
func (p *Project) CreateInitialLists(tx *sql.Tx) error {
	// 初期リストの準備
	closeListOption, err := list_option.FindByAction("close")
	if err != nil {
		tx.Rollback()
		return err
	}
	todo := list.New(
		0,
		p.ID,
		p.UserID,
		config.Element("init_list").(map[interface{}]interface{})["todo"].(string),
		"f37b1d",
		sql.NullInt64{},
		false,
	)
	inprogress := list.New(
		0,
		p.ID,
		p.UserID,
		config.Element("init_list").(map[interface{}]interface{})["inprogress"].(string),
		"5eb95e",
		sql.NullInt64{},
		false,
	)
	done := list.New(
		0,
		p.ID,
		p.UserID,
		config.Element("init_list").(map[interface{}]interface{})["done"].(string),
		"333333",
		sql.NullInt64{Int64: closeListOption.ID, Valid: true},
		false,
	)
	none := list.New(
		0,
		p.ID,
		p.UserID,
		config.Element("init_list").(map[interface{}]interface{})["none"].(string),
		"ffffff",
		sql.NullInt64{},
		false,
	)

	// ここではDBに保存するだけ
	// githubへの同期はこのレイヤーでは行わない
	if err := none.Save(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := todo.Save(tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := inprogress.Save(tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := done.Save(tx); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// DeleteLists delete all lists related a project
func (p *Project) DeleteLists() error {
	lists, err := p.Lists()
	if err != nil {
		return err
	}
	for _, l := range lists {
		err := l.DeleteTasks()
		if err != nil {
			return err
		}
		err = l.Delete()
		if err != nil {
			return err
		}
	}
	noneList, err := p.NoneList()
	err = noneList.DeleteTasks()
	if err != nil {
		return err
	}
	return noneList.Delete()
}

// Delete delete a project model
func (p *Project) Delete() error {
	err := p.infrastructure.Delete()
	if err != nil {
		return err
	}
	p.infrastructure = nil
	return nil
}
