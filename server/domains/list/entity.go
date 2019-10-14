package list

import (
	"database/sql"

	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"

	"github.com/pkg/errors"
)

// List is a entity for list.
type List struct {
	ID        int64
	ProjectID int64
	UserID    int64
	Title     sql.NullString
	Color     sql.NullString
	IsHidden  bool
	Option    *Option
}

// Option has a list option model object
type Option struct {
	ID     int64
	Action string
}

// New returns a new list entity.
func New(id, projectID, userID int64, title, color sql.NullString, isHidden bool, option *Option) *List {
	return &List{
		id,
		projectID,
		userID,
		title,
		color,
		isHidden,
		option,
	}
}

// UpdateExceptInitList update list except initial list
// for example, ToDo, InProgress, and Done
func (l *List) UpdateExceptInitList(title, color sql.NullString, option *Option) error {
	if l.IsInitList() {
		return errors.New("cannot update initial list")
	}

	return l.Update(title, color, option)
}

// Update update list
func (l *List) Update(title, color sql.NullString, option *Option) error {
	if option == nil {
		// It allow list_option to be nullable.
		// When list_option is null, the action doesn't arise.
		logging.SharedInstance().MethodInfo("list", "Update").Debug("cannot find list_options, set null to list_option_id")
	}
	l.Title = title
	l.Color = color
	l.Option = option
	return nil
}

// Hide hide list, and update.
func (l *List) Hide() error {
	l.IsHidden = true
	return nil
}

// Display display list, and update.
func (l *List) Display() error {
	l.IsHidden = false
	return nil
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

// HasCloseAction checks a list has close list option
func (l *List) HasCloseAction() (bool, error) {
	return l.Option.IsCloseAction(), nil
}
