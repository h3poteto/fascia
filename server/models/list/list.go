package list

import (
	"github.com/h3poteto/fascia/server/models/db"

	"database/sql"

	"github.com/pkg/errors"
)

type List struct {
	ID           int64
	ProjectID    int64
	UserID       int64
	Title        sql.NullString
	Color        sql.NullString
	ListOptionID sql.NullInt64
	IsHidden     bool
	database     *sql.DB
}

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

func FindByID(projectID int64, listID int64) (*List, error) {
	database := db.SharedInstance().Connection
	var id, userID int64
	var title, color sql.NullString
	var optionID sql.NullInt64
	var isHidden bool
	rows, err := database.Query("select id, project_id, user_id, title, color, list_option_id, is_hidden from lists where id = ? AND project_id = ?;", listID, projectID)
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

func (u *List) initialize() {
	u.database = db.SharedInstance().Connection
}

func (u *List) Save(tx *sql.Tx) error {
	var err error
	var result sql.Result
	if tx != nil {
		result, err = tx.Exec("insert into lists (project_id, user_id, title, color, list_option_id, is_hidden, created_at) values (?, ?, ?, ?, ?, ?, now());", u.ProjectID, u.UserID, u.Title, u.Color, u.ListOptionID, u.IsHidden)
	} else {
		result, err = u.database.Exec("insert into lists (project_id, user_id, title, color, list_option_id, is_hidden, created_at) values (?, ?, ?, ?, ?, ?, now());", u.ProjectID, u.UserID, u.Title, u.Color, u.ListOptionID, u.IsHidden)
	}
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	u.ID, _ = result.LastInsertId()
	return nil
}

// Update update and save list in database
func (l *List) Update(title, color string, optionID sql.NullInt64) (e error) {
	_, err := l.database.Exec("update lists set title = ?, color = ?, list_option_id = ?, is_hidden = ? where id = ?;", title, color, optionID, l.IsHidden, l.ID)
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
	_, err := l.database.Exec("update lists set is_hidden = true where id = ?;", l.ID)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	l.IsHidden = true
	return nil
}

// Display can display a list, it change is_hidden filed
func (l *List) Display() error {
	_, err := l.database.Exec("update lists set is_hidden = false where id = ?;", l.ID)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	l.IsHidden = false
	return nil
}
