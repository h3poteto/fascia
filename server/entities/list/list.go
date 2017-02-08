package list

import (
	"database/sql"

	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/entities/list_option"
	"github.com/h3poteto/fascia/server/entities/task"
	"github.com/h3poteto/fascia/server/models/db"
	"github.com/h3poteto/fascia/server/models/list"

	"github.com/pkg/errors"
)

// List has a list model object
type List struct {
	ListModel *list.List
	database  *sql.DB
}

// New returns new list entity
func New(id int64, projectID int64, userID int64, title string, color string, optionID sql.NullInt64, isHidden bool) *List {
	return &List{
		ListModel: list.New(id, projectID, userID, title, color, optionID, isHidden),
		database:  db.SharedInstance().Connection,
	}
}

// FindByID returns a list entity
func FindByID(projectID, listID int64) (*List, error) {
	l, err := list.FindByID(projectID, listID)
	if err != nil {
		return nil, err
	}
	return &List{
		ListModel: l,
		database:  db.SharedInstance().Connection,
	}, nil
}

// Save call list model save
func (l *List) Save(tx *sql.Tx) error {
	return l.ListModel.Save(tx)
}

// UpdateExceptInitList update list except initial list
// for example, ToDo, InProgress, and Done
func (l *List) UpdateExceptInitList(title, color string, optionID int64) error {
	// 初期リストに関しては一切編集を許可しない
	// 色は変えられても良いが，titleとactionは変えられては困る
	// 現段階では色も含めてすべて固定とする
	if l.IsInitList() {
		return errors.New("cannot update initial list")
	}

	return l.Update(title, color, optionID)
}

// Update update list
func (l *List) Update(title, color string, optionID int64) error {
	var listOptionID sql.NullInt64
	listOption, err := list_option.FindByID(optionID)
	if err != nil {
		// list_optionはnullでも構わない
		// nullの場合は特にactionが発生しないだけ
		logging.SharedInstance().MethodInfo("list", "Update").Debugf("cannot find list_options, set null to list_option_id: %v", err)
	} else {
		listOptionID = sql.NullInt64{Int64: listOption.ListOptionModel.ID, Valid: true}
	}
	err = l.ListModel.Update(title, color, listOptionID)
	if err != nil {
		return err
	}
	return nil
}

// Hide call list model hide
func (l *List) Hide() error {
	return l.ListModel.Hide()
}

// Display call list model display
func (l *List) Display() error {
	return l.ListModel.Display()
}

// Tasks list up related tasks
func (l *List) Tasks() ([]*task.Task, error) {
	var slice []*task.Task
	rows, err := l.database.Query("select id, list_id, project_id, user_id, issue_number, title, description, pull_request, html_url from tasks where list_id = ? order by display_index;", l.ListModel.ID)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}

	for rows.Next() {
		var id, listID, userID, projectID int64
		var title, description string
		var issueNumber sql.NullInt64
		var pullRequest bool
		var htmlURL sql.NullString
		err := rows.Scan(&id, &listID, &projectID, &userID, &issueNumber, &title, &description, &pullRequest, &htmlURL)
		if err != nil {
			return nil, errors.Wrap(err, "sql select error")
		}
		if listID == l.ListModel.ID {
			l := task.New(id, listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL)
			slice = append(slice, l)
		}
	}
	return slice, nil
}

// ListOption list up a related list option
func (l *List) ListOption() (*list_option.ListOption, error) {
	if !l.ListModel.ListOptionID.Valid {
		return nil, errors.New("list has no list option")
	}
	option, err := list_option.FindByID(l.ListModel.ListOptionID.Int64)
	if err != nil {
		return nil, err
	}
	return option, nil
}

// IsInitList return true when list is initial list
// for example, ToDo, InProgress, and Done
func (l *List) IsInitList() bool {
	for _, elem := range config.Element("init_list").(map[interface{}]interface{}) {
		if l.ListModel.Title.String == elem.(string) {
			return true
		}
	}
	return false
}

// HasCloseAction check a list has close list option
func (l *List) HasCloseAction() (bool, error) {
	option, err := l.ListOption()
	if err != nil {
		return false, err
	}
	return option.IsCloseAction(), nil
}

func (l *List) DeleteTasks() error {
	_, err := l.database.Exec("DELETE FROM tasks WHERE list_id = ?;", l.ListModel.ID)
	if err != nil {
		return err
	}
	return nil
}

func (l *List) Delete() error {
	err := l.ListModel.Delete()
	if err != nil {
		return err
	}
	l.ListModel = nil
	return nil
}
