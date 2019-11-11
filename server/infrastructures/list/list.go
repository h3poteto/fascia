package list

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/server/domains/list"

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
func (l *List) Find(targetProjectID int64, listID int64) (*list.List, error) {
	var id, projectID, userID int64
	var title, color sql.NullString
	var optionID sql.NullInt64
	var isHidden bool
	err := l.db.QueryRow("SELECT id, project_id, user_id, title, color, list_option_id, is_hidden FROM lists WHERE id = $1 AND project_id = $2;", listID, targetProjectID).Scan(&id, &projectID, &userID, &title, &color, &optionID, &isHidden)
	if err != nil {
		return nil, errors.Wrap(err, "list repository")
	}
	if id != listID {
		return nil, errors.New("cannot find list or project did not contain list")
	}

	var option *list.Option
	if optionID.Valid {
		option, err = l.FindOptionByID(optionID.Int64)
		if err != nil {
			return nil, errors.Wrap(err, "list option repository")
		}
	}
	return &list.List{
		ID:        id,
		ProjectID: projectID,
		UserID:    userID,
		Title:     title,
		Color:     color,
		IsHidden:  isHidden,
		Option:    option,
	}, nil
}

// FindByTaskID retruns parent list of a task.
func (l *List) FindByTaskID(taskID int64) (*list.List, error) {
	var id, projectID, userID int64
	var title, color sql.NullString
	var optionID sql.NullInt64
	var isHidden bool
	err := l.db.QueryRow("SELECT lists.id, lists.project_id, lists.user_id, list.title, list.color, list.list_option_id, lists.is_hidden FROM tasks INNER JOIN lists ON tasks.list_id = lists.id WHERE tasks.id = $1;", taskID).Scan(&id, &projectID, &userID, &title, &color, &optionID, &isHidden)
	if err != nil {
		return nil, errors.Wrap(err, "list repository")
	}
	var option *list.Option
	if optionID.Valid {
		option, err = l.FindOptionByID(optionID.Int64)
		if err != nil {
			return nil, errors.Wrap(err, "list option repository")
		}
	}
	return &list.List{
		ID:        id,
		ProjectID: projectID,
		UserID:    userID,
		Title:     title,
		Color:     color,
		IsHidden:  isHidden,
		Option:    option,
	}, nil
}

// Lists returns all lists related a project.
func (l *List) Lists(parentProjectID int64) ([]*list.List, error) {
	rows, err := l.db.Query("SELECT id, project_id, user_id, title, color, list_option_id, is_hidden FROM lists WHERE project_id = $1 AND title != $2;", parentProjectID, config.Element("init_list").(map[interface{}]interface{})["none"].(string))
	if err != nil {
		return nil, errors.Wrap(err, "list repository")
	}
	var result = []*list.List{}
	for rows.Next() {
		var id, projectID, userID int64
		var title, color sql.NullString
		var optionID sql.NullInt64
		var isHidden bool
		err = rows.Scan(&id, &projectID, &userID, &title, &color, &optionID, &isHidden)
		if err != nil {
			return nil, errors.Wrap(err, "list repository")
		}
		// TODO: N+1
		var option *list.Option
		if optionID.Valid {
			option, err = l.FindOptionByID(optionID.Int64)
			if err != nil {
				return nil, errors.Wrap(err, "list option repository")
			}
		}
		if projectID == parentProjectID && title.Valid {
			l := &list.List{
				ID:        id,
				ProjectID: projectID,
				UserID:    userID,
				Title:     title,
				Color:     color,
				IsHidden:  isHidden,
				Option:    option,
			}
			result = append(result, l)
		}
	}
	return result, nil
}

// NoneList returns a none list related a project.
func (l *List) NoneList(parentProjectID int64) (*list.List, error) {
	var id, projectID, userID int64
	var title, color sql.NullString
	var optionID sql.NullInt64
	var isHidden bool
	err := l.db.QueryRow("SELECT id, project_id, user_id, title, color, list_option_id, is_hidden FROM lists WHERE project_id = $1 AND title = $2;", parentProjectID, config.Element("init_list").(map[interface{}]interface{})["none"].(string)).Scan(&id, &projectID, &userID, &title, &color, &optionID, &isHidden)
	if err != nil {
		return nil, errors.Wrap(err, "list repository")
	}
	var option *list.Option
	if optionID.Valid {
		option, err = l.FindOptionByID(optionID.Int64)
		if err != nil {
			return nil, errors.Wrap(err, "list option repository")
		}
	}
	if projectID == parentProjectID && title.Valid {
		return &list.List{
			ID:        id,
			ProjectID: projectID,
			UserID:    userID,
			Title:     title,
			Color:     color,
			IsHidden:  isHidden,
			Option:    option,
		}, nil
	}
	return nil, errors.New("none list not found")
}

// FindOptionByAction search a list option according to action
func (l *List) FindOptionByAction(action string) (*list.Option, error) {
	var id int64
	err := l.db.QueryRow("SELECT id FROM list_options WHERE action = $1;", action).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &list.Option{
		ID:     id,
		Action: action,
	}, nil
}

// FindOptionByID search a list option according to id
func (l *List) FindOptionByID(id int64) (*list.Option, error) {
	var action string
	err := l.db.QueryRow("SELECT action FROM list_options WHERE id = $1;", id).Scan(&action)
	if err != nil {
		return nil, err
	}
	return &list.Option{
		ID:     id,
		Action: action,
	}, nil
}

// AllOption returns all list options.
func (l *List) AllOption() ([]*list.Option, error) {
	rows, err := l.db.Query("SELECT id, action FROM list_options ORDER BY id;")
	if err != nil {
		return nil, err
	}
	result := []*list.Option{}
	for rows.Next() {
		var id int64
		var action string
		err = rows.Scan(&id, &action)
		if err != nil {
			return nil, err
		}
		o := &list.Option{
			ID:     id,
			Action: action,
		}
		result = append(result, o)
	}
	return result, nil
}

// Create save list object to record
func (l *List) Create(projectID int64, userID int64, title sql.NullString, color sql.NullString, listOptionID sql.NullInt64, isHidden bool, tx *sql.Tx) (int64, error) {
	var err error
	var id int64
	if tx != nil {
		err = tx.QueryRow("INSERT INTO lists (project_id, user_id, title, color, list_option_id, is_hidden) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;", projectID, userID, title, color, listOptionID, isHidden).Scan(&id)
	} else {
		err = l.db.QueryRow("INSERT INTO lists (project_id, user_id, title, color, list_option_id, is_hidden) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;", projectID, userID, title, color, listOptionID, isHidden).Scan(&id)
	}
	if err != nil {
		return 0, errors.Wrap(err, "list repository")
	}
	return id, nil
}

// Update update and save list in database
func (l *List) Update(object *list.List) error {
	optionID := sql.NullInt64{}
	if object.Option != nil {
		optionID = sql.NullInt64{
			Int64: object.Option.ID,
			Valid: true,
		}
	}
	_, err := l.db.Exec("UPDATE lists SET project_id = $1, user_id = $2, title = $3, color = $4, list_option_id = $5, is_hidden = $6 WHERE id = $7;", object.ProjectID, object.UserID, object.Title, object.Color, optionID, object.IsHidden, object.ID)

	if err != nil {
		return errors.Wrap(err, "list repository")
	}
	return nil
}

// Delete delete a list model in record
func (l *List) Delete(id int64) error {
	_, err := l.db.Exec("DELETE FROM lists WHERE id = $1;", id)
	if err != nil {
		return errors.Wrap(err, "list repository")
	}
	return nil
}

// DeleteTasks delete all tasks related a list
func (l *List) DeleteTasks(id int64) error {
	_, err := l.db.Exec("DELETE FROM tasks WHERE list_id = $1;", id)
	if err != nil {
		return errors.Wrap(err, "list repository")
	}
	return nil
}
