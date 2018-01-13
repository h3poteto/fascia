package list

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/server/entities/list_option"
	"github.com/h3poteto/fascia/server/entities/task"
	"github.com/h3poteto/fascia/server/infrastructures/list"
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

// Lists returns all lists related a project.
func Lists(projectID int64) ([]*List, error) {
	var slice []*List

	lists, err := list.Lists(projectID)
	if err != nil {
		return nil, err
	}

	for _, l := range lists {
		s := &List{
			infrastructure: l,
		}
		if err := s.reload(); err != nil {
			return nil, err
		}
		slice = append(slice, s)
	}
	return slice, nil
}

// NoneList retruns a none list related a project.
func NoneList(projectID int64) (*List, error) {
	l, err := list.NoneList(projectID)
	if err != nil {
		return nil, err
	}
	s := &List{
		infrastructure: l,
	}
	if err := s.reload(); err != nil {
		return nil, err
	}
	return s, nil
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
