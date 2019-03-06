package list

import (
	"database/sql"

	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"

	"github.com/pkg/errors"
)

// List is a entity for list.
type List struct {
	ID             int64
	ProjectID      int64
	UserID         int64
	Title          sql.NullString
	Color          sql.NullString
	ListOptionID   sql.NullInt64
	IsHidden       bool
	infrastructure Repository
}

type Repository interface {
	Find(int64, int64) (int64, int64, int64, sql.NullString, sql.NullString, sql.NullInt64, bool, error)
	FindByTaskID(int64) (int64, int64, int64, sql.NullString, sql.NullString, sql.NullInt64, bool, error)
	Create(int64, int64, sql.NullString, sql.NullString, sql.NullInt64, bool, *sql.Tx) (int64, error)
	Update(int64, int64, int64, sql.NullString, sql.NullString, sql.NullInt64, bool) error
	Delete(int64) error
	DeleteTasks(int64) error
	Lists(int64) ([]map[string]interface{}, error)
	NoneList(int64) (int64, int64, int64, sql.NullString, sql.NullString, sql.NullInt64, bool, error)
	FindOptionByAction(string) (int64, string, error)
	FindOptionByID(int64) (int64, string, error)
	AllOption() ([]map[string]interface{}, error)
}

// New returns a new list entity.
func New(id, projectID, userID int64, title, color sql.NullString, optionID sql.NullInt64, isHidden bool, infrastructure Repository) *List {
	return &List{
		id,
		projectID,
		userID,
		title,
		color,
		optionID,
		isHidden,
		infrastructure,
	}
}

// Create call list model save
func (l *List) Create(tx *sql.Tx) error {
	id, err := l.infrastructure.Create(l.ProjectID, l.UserID, l.Title, l.Color, l.ListOptionID, l.IsHidden, tx)
	if err != nil {
		return err
	}
	l.ID = id
	return nil
}

// UpdateExceptInitList update list except initial list
// for example, ToDo, InProgress, and Done
func (l *List) UpdateExceptInitList(title, color sql.NullString, optionID int64) error {
	if l.IsInitList() {
		return errors.New("cannot update initial list")
	}

	return l.Update(title, color, optionID)
}

// Update update list
func (l *List) Update(title, color sql.NullString, optionID int64) error {
	var listOptionID sql.NullInt64
	listOption, err := FindOptionByID(optionID, l.infrastructure)
	if err != nil {
		// It allow list_option to be nullable.
		// When list_option is null, the action doesn't arise.
		logging.SharedInstance().MethodInfo("list", "Update").Debugf("cannot find list_options, set null to list_option_id: %v", err)
	} else {
		listOptionID = sql.NullInt64{Int64: listOption.ID, Valid: true}
	}
	err = l.infrastructure.Update(l.ID, l.ProjectID, l.UserID, title, color, listOptionID, l.IsHidden)
	if err != nil {
		return err
	}
	l.Title = title
	l.Color = color
	l.ListOptionID = listOptionID
	return nil
}

// Hide call list model hide
func (l *List) Hide() error {
	err := l.infrastructure.Update(l.ID, l.ProjectID, l.UserID, l.Title, l.Color, l.ListOptionID, true)
	if err != nil {
		return err
	}
	l.IsHidden = true
	return nil
}

// Display call list model display
func (l *List) Display() error {
	err := l.infrastructure.Update(l.ID, l.ProjectID, l.UserID, l.Title, l.Color, l.ListOptionID, false)
	if err != nil {
		return err
	}
	l.IsHidden = false
	return nil
}

// DeleteTasks delete all tasks related a list
func (l *List) DeleteTasks() error {
	return l.infrastructure.DeleteTasks(l.ID)
}

// Delete delete a list model
func (l *List) Delete() error {
	return l.infrastructure.Delete(l.ID)
}

// IsInitList return true when list is initial list
// for example, ToDo, InProgress, and Done
func (l *List) IsInitList() bool {
	for _, elem := range config.Element("init_list").(map[interface{}]interface{}) {
		if l.Title.String == elem.(string) {
			return true
		}
	}
	return false
}
