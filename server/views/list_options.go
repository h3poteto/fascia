package views

import (
	"github.com/h3poteto/fascia/server/entities/list_option"
)

type ListOption struct {
	ID     int64  `json:ID`
	Action string `json:Action`
}

func ParseListOptionJson(option *list_option.ListOption) (*ListOption, error) {
	return &ListOption{
		ID:     option.ListOptionModel.ID,
		Action: option.ListOptionModel.Action,
	}, nil
}

func ParseListOptionsJson(options []*list_option.ListOption) ([]*ListOption, error) {
	results := make([]*ListOption, 0)
	for _, o := range options {
		parse, err := ParseListOptionJson(o)
		if err != nil {
			return nil, err
		}
		results = append(results, parse)
	}
	return results, nil
}
