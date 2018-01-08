package list

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/server/entities/list_option"
	"github.com/h3poteto/fascia/server/entities/task"
	"github.com/pkg/errors"
)

// FindByID returns a list entity
func FindByID(projectID, listID int64) (*List, error) {
	l := &List{
		ID:        listID,
		ProjectID: projectID,
	}
	err := l.reload()
	if err != nil {
		return nil, err
	}
	return l, nil
}

// Tasks list up related tasks
func (l *List) Tasks() ([]*task.Task, error) {
	return task.Tasks(l.ID)
}

// ListOption list up a related list option
func (l *List) ListOption() (*list_option.ListOption, error) {
	if !l.infrastructure.ListOptionID.Valid {
		return nil, errors.New("list has no list option")
	}
	option, err := list_option.FindByID(l.infrastructure.ListOptionID.Int64)
	if err != nil {
		return nil, err
	}
	return option, nil
}

// IsInitList return true when list is initial list
// for example, ToDo, InProgress, and Done
func (l *List) IsInitList() bool {
	for _, elem := range config.Element("init_list").(map[interface{}]interface{}) {
		if l.infrastructure.Title.String == elem.(string) {
			return true
		}
	}
	return false
}

// HasCloseAction check a list has close list option
func (l *List) HasCloseAction() (bool, error) {
	option, err := l.ListOption()
	if err != nil {
		return false, err
	}
	return option.IsCloseAction(), nil
}
