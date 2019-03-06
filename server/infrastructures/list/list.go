package list

import (
	"github.com/h3poteto/fascia/config"

	"database/sql"

	"github.com/pkg/errors"
)

// List has list record
type List struct {
	db *sql.DB
}

// New returns a new list object
func New(db *sql.DB) *List {
	return &List{
		db,
	}
}

// Find search a list according to id
func (l *List) Find(targetProjectID int64, listID int64) (int64, int64, int64, sql.NullString, sql.NullString, sql.NullInt64, bool, error) {
	var id, projectID, userID int64
	var title, color sql.NullString
	var optionID sql.NullInt64
	var isHidden bool
	err := l.db.QueryRow("select id, project_id, user_id, title, color, list_option_id, is_hidden from lists where id = ? AND project_id = ?;", listID, targetProjectID).Scan(&id, &projectID, &userID, &title, &color, &optionID, &isHidden)
	if err != nil {
		return 0, 0, 0, sql.NullString{}, sql.NullString{}, sql.NullInt64{}, false, errors.Wrap(err, "list repository")
	}
	if id != listID {
		return 0, 0, 0, sql.NullString{}, sql.NullString{}, sql.NullInt64{}, false, errors.New("cannot find list or project did not contain list")
	}
	return id, projectID, userID, title, color, optionID, isHidden, nil
}

// FindByTaskID retruns parent list of a task.
func (l *List) FindByTaskID(taskID int64) (int64, int64, int64, sql.NullString, sql.NullString, sql.NullInt64, bool, error) {
	var id, projectID, userID int64
	var title, color sql.NullString
	var optionID sql.NullInt64
	var isHidden bool
	err := l.db.QueryRow("SELECT lists.id, lists.project_id, lists.user_id, list.title, list.color, list.list_option_id, lists.is_hidden FROM tasks INNER JOIN lists on tasks.list_id = lists.id WHERE tasks.id = ?;", taskID).Scan(&id, &projectID, &userID, &title, &color, &optionID, &isHidden)
	if err != nil {
		return 0, 0, 0, sql.NullString{}, sql.NullString{}, sql.NullInt64{}, false, errors.Wrap(err, "list repository")
	}
	return id, projectID, userID, title, color, optionID, isHidden, nil
}

// Create save list object to record
func (l *List) Create(projectID int64, userID int64, title sql.NullString, color sql.NullString, listOptionID sql.NullInt64, isHidden bool, tx *sql.Tx) (int64, error) {
	var err error
	var result sql.Result
	if tx != nil {
		result, err = tx.Exec("insert into lists (project_id, user_id, title, color, list_option_id, is_hidden, created_at) values (?, ?, ?, ?, ?, ?, now());", projectID, userID, title, color, listOptionID, isHidden)
	} else {
		result, err = l.db.Exec("insert into lists (project_id, user_id, title, color, list_option_id, is_hidden, created_at) values (?, ?, ?, ?, ?, ?, now());", projectID, userID, title, color, listOptionID, isHidden)
	}
	if err != nil {
		return 0, errors.Wrap(err, "list repository")
	}
	id, _ := result.LastInsertId()
	return id, nil
}

// Update update and save list in database
func (l *List) Update(id int64, projectID int64, userID int64, title sql.NullString, color sql.NullString, listOptionID sql.NullInt64, isHidden bool) error {
	_, err := l.db.Exec("update lists set project_id = ?, user_id = ?, title = ?, color = ?, list_option_id = ?, is_hidden = ? where id = ?;", projectID, userID, title, color, listOptionID, isHidden, id)
	if err != nil {
		return errors.Wrap(err, "list repository")
	}
	return nil
}

// Delete delete a list model in record
func (l *List) Delete(id int64) error {
	_, err := l.db.Exec("DELETE FROM lists WHERE id = ?;", id)
	if err != nil {
		return errors.Wrap(err, "list repository")
	}
	return nil
}

// DeleteTasks delete all tasks related a list
func (l *List) DeleteTasks(id int64) error {
	_, err := l.db.Exec("DELETE FROM tasks WHERE list_id = ?;", id)
	if err != nil {
		return errors.Wrap(err, "list repository")
	}
	return nil
}

// Lists returns all lists related a project.
func (l *List) Lists(parentProjectID int64) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	rows, err := l.db.Query("select id, project_id, user_id, title, color, list_option_id, is_hidden from lists where project_id = ? and title != ?;", parentProjectID, config.Element("init_list").(map[interface{}]interface{})["none"].(string))
	if err != nil {
		return result, errors.Wrap(err, "list repository")
	}
	for rows.Next() {
		var id, projectID, userID int64
		var title, color sql.NullString
		var optionID sql.NullInt64
		var isHidden bool
		err = rows.Scan(&id, &projectID, &userID, &title, &color, &optionID, &isHidden)
		if err != nil {
			return nil, errors.Wrap(err, "list repository")
		}
		if projectID == parentProjectID && title.Valid {
			l := map[string]interface{}{
				"id":        id,
				"projectID": projectID,
				"userID":    userID,
				"title":     title,
				"color":     color,
				"optionID":  optionID,
				"isHidden":  isHidden,
			}
			result = append(result, l)
		}
	}
	return result, nil
}

// NoneList returns a none list related a project.
func (l *List) NoneList(parentProjectID int64) (int64, int64, int64, sql.NullString, sql.NullString, sql.NullInt64, bool, error) {
	var id, projectID, userID int64
	var title, color sql.NullString
	var optionID sql.NullInt64
	var isHidden bool
	err := l.db.QueryRow("select id, project_id, user_id, title, color, list_option_id, is_hidden from lists where project_id = ? and title = ?;", parentProjectID, config.Element("init_list").(map[interface{}]interface{})["none"].(string)).Scan(&id, &projectID, &userID, &title, &color, &optionID, &isHidden)
	if err != nil {
		return 0, 0, 0, sql.NullString{}, sql.NullString{}, sql.NullInt64{}, false, err
	}
	if projectID == parentProjectID && title.Valid {
		return id, projectID, userID, title, color, optionID, isHidden, nil
	}
	return 0, 0, 0, sql.NullString{}, sql.NullString{}, sql.NullInt64{}, false, errors.New("none list not found")
}

// FindOptionByAction search a list option according to action
func (l *List) FindOptionByAction(action string) (int64, string, error) {
	var id int64
	err := l.db.QueryRow("select id from list_options where action = ?;", action).Scan(&id)
	if err != nil {
		return 0, "", err
	}
	return id, action, nil
}

// FindOptionByID search a list option according to id
func (l *List) FindOptionByID(id int64) (int64, string, error) {
	var action string
	err := l.db.QueryRow("select action from list_options where id = ?;", id).Scan(&action)
	if err != nil {
		return 0, "", err
	}
	return id, action, nil
}

// AllOption returns all list options.
func (l *List) AllOption() ([]map[string]interface{}, error) {
	slice := []map[string]interface{}{}
	rows, err := l.db.Query("select id, action from list_options;")
	if err != nil {
		return slice, err
	}
	for rows.Next() {
		var id int64
		var action string
		err = rows.Scan(&id, &action)
		if err != nil {
			return nil, err
		}
		m := map[string]interface{}{"id": id, "action": action}
		slice = append(slice, m)
	}
	return slice, nil
}
