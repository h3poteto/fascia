package board

import (
	"github.com/h3poteto/fascia/server/domains/entities/list_option"
)

// ListOption has a list option entity
type ListOption struct {
	ListOptionEntity *list_option.ListOption
}

// ListOptionAll returns all list options
func ListOptionAll() ([]*ListOption, error) {
	options, err := list_option.ListOptionAll()
	if err != nil {
		return nil, err
	}
	var listOptions []*ListOption
	for _, o := range options {
		listOptions = append(listOptions, &ListOption{ListOptionEntity: o})
	}
	return listOptions, nil
}

// FindListOptionByID returns a list option service
func FindListOptionByID(id int64) (*ListOption, error) {
	option, err := list_option.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &ListOption{
		ListOptionEntity: option,
	}, nil
}

// FindListOptionByAction returns a list option service
func FindListOptionByAction(action string) (*ListOption, error) {
	option, err := list_option.FindByAction(action)
	if err != nil {
		return nil, err
	}
	return &ListOption{
		ListOptionEntity: option,
	}, nil
}

// IsCloseAction returns either close action or other action
func (l *ListOption) IsCloseAction() bool {
	return l.ListOptionEntity.IsCloseAction()
}
