package services

import (
	"github.com/h3poteto/fascia/server/aggregations/list_option"
)

type ListOption struct {
	ListOptionAggregation *list_option.ListOption
}

func ListOptionAll() ([]*ListOption, error) {
	options, err := list_option.ListOptionAll()
	if err != nil {
		return nil, err
	}
	var listOptions []*ListOption
	for _, o := range options {
		listOptions = append(listOptions, &ListOption{ListOptionAggregation: o})
	}
	return listOptions, nil
}

func FindListOptionByID(id int64) (*ListOption, error) {
	option, err := list_option.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &ListOption{
		ListOptionAggregation: option,
	}, nil
}

func FindListOptionByAction(action string) (*ListOption, error) {
	option, err := list_option.FindByAction(action)
	if err != nil {
		return nil, err
	}
	return &ListOption{
		ListOptionAggregation: option,
	}, nil
}

func (l *ListOption) IsCloseAction() bool {
	return l.ListOptionAggregation.IsCloseAction()
}
