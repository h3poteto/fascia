package views

import (
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/usecases/board"
)

// List provides a response structure for list
type List struct {
	ID           int64   `json:ID`
	ProjectID    int64   `json:ProjectID`
	UserID       int64   `json:UserID`
	Title        string  `json:Title`
	ListTasks    []*Task `json:ListTasks`
	Color        string  `json:Color`
	ListOptionID int64   `json:ListOptionID`
	IsHidden     bool    `json:IsHidden`
	IsInitList   bool    `json:IsInitList`
}

// AllLists providers a response structure for lists and none lists
type AllLists struct {
	Lists    []*List `json:Lists`
	NoneList *List   `json:NoneList`
}

// ParseListJSON returns a List struct for response
func ParseListJSON(list *list.List) (*List, error) {
	tasks, err := board.ListTasks(list)
	if err != nil {
		return nil, err
	}
	jsonTasks, err := ParseTasksJSON(tasks)
	if err != nil {
		return nil, err
	}
	if list.Option != nil {
		return &List{
			ID:           list.ID,
			ProjectID:    list.ProjectID,
			UserID:       list.UserID,
			Title:        list.Title.String,
			ListTasks:    jsonTasks,
			Color:        list.Color.String,
			ListOptionID: list.Option.ID,
			IsHidden:     list.IsHidden,
			IsInitList:   list.IsInitList(),
		}, nil
	}
	return &List{
		ID:           list.ID,
		ProjectID:    list.ProjectID,
		UserID:       list.UserID,
		Title:        list.Title.String,
		ListTasks:    jsonTasks,
		Color:        list.Color.String,
		ListOptionID: 0,
		IsHidden:     list.IsHidden,
		IsInitList:   list.IsInitList(),
	}, nil
}

// ParseListsJSON returns some List structs for response
func ParseListsJSON(lists []*list.List) ([]*List, error) {
	results := []*List{}
	for _, l := range lists {
		parse, err := ParseListJSON(l)
		if err != nil {
			return nil, err
		}
		results = append(results, parse)
	}
	return results, nil
}

// ParseAllListsJSON returns a AllLists struct for response
func ParseAllListsJSON(noneList *list.List, lists []*list.List) (*AllLists, error) {
	jsonNone, err := ParseListJSON(noneList)
	if err != nil {
		return nil, err
	}
	jsonLists, err := ParseListsJSON(lists)
	if err != nil {
		return nil, err
	}
	return &AllLists{Lists: jsonLists, NoneList: jsonNone}, nil
}
