package repositories

import (
	"database/sql"
	"errors"
)

type ListOption struct {
	ID     int64
	Action string
}

type DummyList struct {
	ProjectID int64
	UserID    int64
	Option    *ListOption
}

func (d *DummyList) Find(projectID, id int64) (int64, int64, int64, sql.NullString, sql.NullString, sql.NullInt64, bool, error) {
	return id, projectID, d.UserID, sql.NullString{String: "title", Valid: true}, sql.NullString{String: "#ff0000", Valid: true}, sql.NullInt64{Int64: d.Option.ID, Valid: true}, false, nil
}

func (d *DummyList) FindByTaskID(taskID int64) (int64, int64, int64, sql.NullString, sql.NullString, sql.NullInt64, bool, error) {
	return 1, d.ProjectID, d.UserID, sql.NullString{String: "title", Valid: true}, sql.NullString{String: "#ff0000", Valid: true}, sql.NullInt64{Int64: d.Option.ID, Valid: true}, false, nil
}

func (d *DummyList) Create(projectID, userID int64, title, color sql.NullString, listOptionID sql.NullInt64, isHidden bool, tx *sql.Tx) (int64, error) {
	return 1, nil
}

func (d *DummyList) Update(id, projectID, userID int64, title, color sql.NullString, listOptionID sql.NullInt64, isHidden bool) error {
	return nil
}

func (d *DummyList) Delete(id int64) error {
	return nil
}

func (d *DummyList) DeleteTasks(id int64) error {
	return nil
}

func (d *DummyList) Lists(projectID int64) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	l := map[string]interface{}{
		"id":        1,
		"projectID": projectID,
		"userID":    d.UserID,
		"title":     sql.NullString{String: "title", Valid: true},
		"color":     sql.NullString{String: "#ff0000", Valid: true},
		"optionID":  d.Option.ID,
		"isHidden":  false,
	}
	result = append(result, l)
	return result, nil
}

func (d *DummyList) NoneList(projectID int64) (int64, int64, int64, sql.NullString, sql.NullString, sql.NullInt64, bool, error) {
	return 1, d.ProjectID, d.UserID, sql.NullString{String: "None", Valid: true}, sql.NullString{String: "#ff0000", Valid: true}, sql.NullInt64{Int64: d.Option.ID, Valid: true}, false, nil
}

func (d *DummyList) FindOptionByAction(action string) (int64, string, error) {
	if d.Option.Action != action {
		return 0, "", errors.New("option does not exist")
	}
	return d.Option.ID, action, nil
}

func (d *DummyList) FindOptionByID(id int64) (int64, string, error) {
	if d.Option.ID != id {
		return 0, "", errors.New("option does not exist")
	}
	return id, d.Option.Action, nil
}

func (d *DummyList) AllOption() ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	o := map[string]interface{}{
		"id":     1,
		"action": "TODO",
	}
	result = append(result, o)
	return result, nil
}
