package views

import "github.com/h3poteto/fascia/server/domains/list"

// ListOption provides a response structure for list option
type ListOption struct {
	ID     int64  `json:ID`
	Action string `json:Action`
}

// ParseListOptionJSON returns a ListOption struct for response
func ParseListOptionJSON(option *list.Option) (*ListOption, error) {
	return &ListOption{
		ID:     option.ID,
		Action: option.Action,
	}, nil
}

// ParseListOptionsJSON returns some ListOption structs for response
func ParseListOptionsJSON(options []*list.Option) ([]*ListOption, error) {
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
