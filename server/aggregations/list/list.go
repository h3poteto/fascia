package list

import (
	"database/sql"

	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/aggregations/list_option"
	"github.com/h3poteto/fascia/server/aggregations/task"
	"github.com/h3poteto/fascia/server/models/db"
	"github.com/h3poteto/fascia/server/models/list"

	"github.com/pkg/errors"
)

type List struct {
	ListModel *list.List
	database  *sql.DB
}

func New(id int64, projectID int64, userID int64, title string, color string, optionID sql.NullInt64, isHidden bool) *List {
	return &List{
		ListModel: list.New(id, projectID, userID, title, color, optionID, isHidden),
		database:  db.SharedInstance().Connection,
	}
}

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

func (l *List) Save(tx *sql.Tx) error {
	return l.ListModel.Save(tx)
}

func (l *List) Update(title, color string, optionID int64) error {
	// 初期リストに関しては一切編集を許可しない
	// 色は変えられても良いが，titleとactionは変えられては困る
	// 現段階では色も含めてすべて固定とする
	if l.IsInitList() {
		return errors.New("cannot update initial list")
	}

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

func (l *List) Hide() error {
	return l.ListModel.Hide()
}

func (l *List) Display() error {
	return l.ListModel.Display()
}

// Tasks list up related a list
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

// ListOption
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

func (l *List) IsInitList() bool {
	for _, elem := range config.Element("init_list").(map[interface{}]interface{}) {
		if l.ListModel.Title.String == elem.(string) {
			return true
		}
	}
	return false
}

// HasCloseAction check a list has close list_option
func (l *List) HasCloseAction() (bool, error) {
	option, err := l.ListOption()
	if err != nil {
		return false, err
	}
	return option.IsCloseAction(), nil
}
