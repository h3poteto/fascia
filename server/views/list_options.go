package views

import (
	"github.com/h3poteto/fascia/server/entities/list_option"
)

// ListOption provides a response structure for list option
type ListOption struct {
	ID     int64  `json:ID`
	Action string `json:Action`
}

// ParseListOptionJSON returns a ListOption struct for response
func ParseListOptionJSON(option *list_option.ListOption) (*ListOption, error) {
	return &ListOption{
		ID:     option.ListOptionModel.ID,
		Action: option.ListOptionModel.Action,
	}, nil
}

// ParseListOptionsJSON returns some ListOption structs for response
func ParseListOptionsJSON(options []*list_option.ListOption) ([]*ListOption, error) {
	results := []*ListOption{}
	for _, o := range options {
		parse, err := ParseListOptionJSON(o)
		if err != nil {
			return nil, err
		}
		results = append(results, parse)
	}
	return results, nil
}
