package views

import (
	"github.com/h3poteto/fascia/server/entities/list"
)

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

type AllLists struct {
	Lists    []*List `json:Lists`
	NoneList *List   `json:NoneList`
}

func ParseListJSON(list *list.List) (*List, error) {
	tasks, err := list.Tasks()
	if err != nil {
		return nil, err
	}
	jsonTasks, err := ParseTasksJSON(tasks)
	if err != nil {
		return nil, err
	}
	return &List{
		ID:           list.ListModel.ID,
		ProjectID:    list.ListModel.ProjectID,
		UserID:       list.ListModel.UserID,
		Title:        list.ListModel.Title.String,
		ListTasks:    jsonTasks,
		Color:        list.ListModel.Color.String,
		ListOptionID: list.ListModel.ListOptionID.Int64,
		IsHidden:     list.ListModel.IsHidden,
		IsInitList:   list.IsInitList(),
	}, nil
}

func ParseListsJSON(lists []*list.List) ([]*List, error) {
	var results []*List
	for _, l := range lists {
		parse, err := ParseListJSON(l)
		if err != nil {
			return nil, err
		}
		results = append(results, parse)
	}
	return results, nil
}

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
