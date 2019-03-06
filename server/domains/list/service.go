package list

import (
	"database/sql"

	"github.com/h3poteto/fascia/server/domains/task"
	"github.com/pkg/errors"
)

// Find returns a list entity
func Find(targetProjectID, targetListID int64, infrastructure Repository) (*List, error) {
	id, projectID, userID, title, color, optionID, isHidden, err := infrastructure.Find(targetProjectID, targetListID)
	if err != nil {
		return nil, err
	}
	return New(id, projectID, userID, title, color, optionID, isHidden, infrastructure), nil
}

func FindByTaskID(targetTaskID int64, infrastructure Repository) (*List, error) {
	id, projectID, userID, title, color, optionID, isHidden, err := infrastructure.FindByTaskID(targetTaskID)
	if err != nil {
		return nil, err
	}
	return New(id, projectID, userID, title, color, optionID, isHidden, infrastructure), nil
}

// Lists returns all lists related a project.
func Lists(targetProjectID int64, infrastructure Repository) ([]*List, error) {
	var result []*List

	lists, err := infrastructure.Lists(targetProjectID)
	if err != nil {
		return result, err
	}
	for _, list := range lists {
		id, ok := list["id"].(int64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		projectID, ok := list["projectID"].(int64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		userID, ok := list["userID"].(int64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		title, ok := list["title"].(sql.NullString)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		color, ok := list["color"].(sql.NullString)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		optionID, ok := list["optionID"].(sql.NullInt64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		isHidden, ok := list["isHidden"].(bool)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		l := New(id, projectID, userID, title, color, optionID, isHidden, infrastructure)
		result = append(result, l)
	}
	return result, nil
}

// NoneList retruns a none list related a project.
func NoneList(targetProjectID int64, infrastructure Repository) (*List, error) {
	id, projectID, userID, title, color, optionID, isHidden, err := infrastructure.NoneList(targetProjectID)
	if err != nil {
		return nil, err
	}
	return New(id, projectID, userID, title, color, optionID, isHidden, infrastructure), nil
}

func (l *List) ListOption() (*Option, error) {
	if !l.ListOptionID.Valid {
		return nil, errors.New("list has no list option")
	}
	option, err := FindOptionByID(l.ListOptionID.Int64, l.infrastructure)
	if err != nil {
		return nil, err
	}
	return option, nil
}

// // HasCloseAction check a list has close list option
func (l *List) HasCloseAction() (bool, error) {
	option, err := l.ListOption()
	if err != nil {
		return false, err
	}
	return option.IsCloseAction(), nil
}

func (l *List) Tasks(infrastructure task.Repository) ([]*task.Task, error) {
	return task.Tasks(l.ID, infrastructure)
}
