package board

import (
	"database/sql"

	"github.com/h3poteto/fascia/lib/modules/database"
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/domains/services"
	"github.com/h3poteto/fascia/server/domains/task"
	repository "github.com/h3poteto/fascia/server/infrastructures/list"
)

// InjectDB set DB connection from connection pool.
func InjectDB() *sql.DB {
	return database.SharedInstance().Connection
}

// InjectListRepository returns list Repository.
func InjectListRepository() list.Repository {
	return repository.New(InjectDB())
}

// ListHasCloseAction returns either close action or not.
func ListHasCloseAction(list *list.List) (bool, error) {
	return list.HasCloseAction()
}

// FindList returns a list.
func FindList(projectID, listID int64) (*list.List, error) {
	repo := InjectListRepository()
	return repo.Find(projectID, listID)
}

// CreateList creates a list, and sync to github.
func CreateList(projectID, userID int64, title, color string, optionID sql.NullInt64, isHidden bool) (*list.List, error) {
	nullableTitle := sql.NullString{String: title, Valid: true}
	nullableColor := sql.NullString{String: color, Valid: true}
	repo := InjectListRepository()
	id, err := repo.Create(projectID, userID, nullableTitle, nullableColor, optionID, isHidden, nil)
	if err != nil {
		return nil, err
	}
	l, err := repo.Find(projectID, id)
	if err != nil {
		return nil, err
	}

	go services.AfterCreateList(l, InjectProjectRepository(), InjectRepoRepository())
	return l, nil
}

// UpdateList updates a list, and sync to github.
func UpdateList(l *list.List, title, color string, optionID int64) (*list.List, error) {
	nullableTitle := sql.NullString{String: title, Valid: true}
	nullableColor := sql.NullString{String: color, Valid: true}
	repo := InjectListRepository()
	// Allow null option, so ignore error.
	option, _ := repo.FindOptionByID(optionID)
	err := l.Update(nullableTitle, nullableColor, option)
	if err != nil {
		return nil, err
	}
	if err = repo.Update(l); err != nil {
		return nil, err
	}

	go services.AfterUpdateList(l, title, color, InjectProjectRepository(), InjectRepoRepository())
	return l, nil
}

// ListTasks returns tasks related the list.
func ListTasks(l *list.List) ([]*task.Task, error) {
	repo := InjectTaskRepository()
	return repo.Tasks(l.ID)
}

// ListOptionAll returns all list options
func ListOptionAll() ([]*list.Option, error) {
	repo := InjectListRepository()
	return repo.AllOption()
}

// FindListOptionByID returns a list option service
func FindListOptionByID(id int64) (*list.Option, error) {
	repo := InjectListRepository()
	return repo.FindOptionByID(id)
}

// FindListOptionByAction returns a list option service
func FindListOptionByAction(action string) (*list.Option, error) {
	repo := InjectListRepository()
	return repo.FindOptionByAction(action)
}

// HideList hides a list.
func HideList(l *list.List) error {
	l.Hide()
	repo := InjectListRepository()
	return repo.Update(l)
}

// DisplayList display a list.
func DisplayList(l *list.List) error {
	l.Display()
	repo := InjectListRepository()
	return repo.Update(l)
}
