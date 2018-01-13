package list

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/database"

	"database/sql"

	"github.com/pkg/errors"
)

// List has list record
type List struct {
	ID           int64
	ProjectID    int64
	UserID       int64
	Title        sql.NullString
	Color        sql.NullString
	ListOptionID sql.NullInt64
	IsHidden     bool
	db           *sql.DB
}

// New returns a new list object
func New(id int64, projectID int64, userID int64, title string, color string, optionID sql.NullInt64, isHidden bool) *List {
	if projectID == 0 {
		return nil
	}
	nullTitle := sql.NullString{String: title, Valid: true}
	nullColor := sql.NullString{String: color, Valid: true}

	list := &List{ID: id, ProjectID: projectID, UserID: userID, Title: nullTitle, Color: nullColor, ListOptionID: optionID, IsHidden: isHidden}
	list.initialize()
	return list
}

// FindByID search a list according to id
func FindByID(projectID int64, listID int64) (*List, error) {
	db := database.SharedInstance().Connection
	var id, userID int64
	var title, color sql.NullString
	var optionID sql.NullInt64
	var isHidden bool
	rows, err := db.Query("select id, project_id, user_id, title, color, list_option_id, is_hidden from lists where id = ? AND project_id = ?;", listID, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	for rows.Next() {
		err = rows.Scan(&id, &projectID, &userID, &title, &color, &optionID, &isHidden)
		if err != nil {
			return nil, errors.Wrap(err, "sql scan error")
		}
	}
	if id != listID {
		return nil, errors.New("cannot find list or project did not contain list")
	}
	list := New(id, projectID, userID, title.String, color.String, optionID, isHidden)
	return list, nil

}

func (l *List) initialize() {
	l.db = database.SharedInstance().Connection
}

// Save save list object to record
func (l *List) Save(tx *sql.Tx) error {
	var err error
	var result sql.Result
	if tx != nil {
		result, err = tx.Exec("insert into lists (project_id, user_id, title, color, list_option_id, is_hidden, created_at) values (?, ?, ?, ?, ?, ?, now());", l.ProjectID, l.UserID, l.Title, l.Color, l.ListOptionID, l.IsHidden)
	} else {
		result, err = l.db.Exec("insert into lists (project_id, user_id, title, color, list_option_id, is_hidden, created_at) values (?, ?, ?, ?, ?, ?, now());", l.ProjectID, l.UserID, l.Title, l.Color, l.ListOptionID, l.IsHidden)
	}
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	l.ID, _ = result.LastInsertId()
	return nil
}

// Update update and save list in database
func (l *List) Update(title, color string, optionID sql.NullInt64) (e error) {
	_, err := l.db.Exec("update lists set title = ?, color = ?, list_option_id = ?, is_hidden = ? where id = ?;", title, color, optionID, l.IsHidden, l.ID)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}

	l.Title = sql.NullString{String: title, Valid: true}
	l.Color = sql.NullString{String: color, Valid: true}
	l.ListOptionID = optionID
	return nil
}

// Hide can hide a list, it change is_hidden field
func (l *List) Hide() error {
	_, err := l.db.Exec("update lists set is_hidden = true where id = ?;", l.ID)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	l.IsHidden = true
	return nil
}

// Display can display a list, it change is_hidden filed
func (l *List) Display() error {
	_, err := l.db.Exec("update lists set is_hidden = false where id = ?;", l.ID)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	l.IsHidden = false
	return nil
}

// Delete delete a list model in record
func (l *List) Delete() error {
	_, err := l.db.Exec("DELETE FROM lists WHERE id = ?;", l.ID)
	if err != nil {
		return errors.Wrap(err, "list delete error")
	}
	return nil
}

// DeleteTasks delete all tasks related a list
func (l *List) DeleteTasks() error {
	_, err := l.db.Exec("DELETE FROM tasks WHERE list_id = ?;", l.ID)
	if err != nil {
		return err
	}
	return nil
}

// Lists returns all lists related a project.
func Lists(parentProjectID int64) ([]*List, error) {
	db := database.SharedInstance().Connection
	var slice []*List
	rows, err := db.Query("select id, project_id, user_id, title, color, list_option_id, is_hidden from lists where project_id = ? and title != ?;", parentProjectID, config.Element("init_list").(map[interface{}]interface{})["none"].(string))
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	for rows.Next() {
		var id, projectID, userID int64
		var title, color sql.NullString
		var optionID sql.NullInt64
		var isHidden bool
		err = rows.Scan(&id, &projectID, &userID, &title, &color, &optionID, &isHidden)
		if err != nil {
			return nil, errors.Wrap(err, "sql select error")
		}
		if projectID == parentProjectID && title.Valid {
			l := New(id, projectID, userID, title.String, color.String, optionID, isHidden)
			slice = append(slice, l)
		}
	}
	return slice, nil
}

// NoneList returns a none list related a project.
func NoneList(parentProjectID int64) (*List, error) {
	db := database.SharedInstance().Connection
	var id, projectID, userID int64
	var title, color sql.NullString
	var optionID sql.NullInt64
	var isHidden bool
	err := db.QueryRow("select id, project_id, user_id, title, color, list_option_id, is_hidden from lists where project_id = ? and title = ?;", parentProjectID, config.Element("init_list").(map[interface{}]interface{})["none"].(string)).Scan(&id, &projectID, &userID, &title, &color, &optionID, &isHidden)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	if projectID == parentProjectID && title.Valid {
		return New(id, projectID, userID, title.String, color.String, optionID, isHidden), nil
	}
	return nil, errors.New("none list not found")
}
