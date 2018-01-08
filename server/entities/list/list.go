package list

import (
	"database/sql"

	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/entities/list_option"
	"github.com/h3poteto/fascia/server/infrastructures/list"

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
	infrastructure *list.List
}

// New returns a new list entity.
// TODO: idなしの新規オブジェクト生成のためだけの関数にしたい
func New(id, projectID, userID int64, title, color string, optionID sql.NullInt64, isHidden bool) *List {
	infrastructure := list.New(id, projectID, userID, title, color, optionID, isHidden)
	l := &List{
		infrastructure: infrastructure,
	}
	l.reload()
	return l
}

// reflect the latest state in infrastructure.
// It is a mapping function.
func (l *List) reflect() {
	l.infrastructure.ID = l.ID
	l.infrastructure.ProjectID = l.ProjectID
	l.infrastructure.UserID = l.UserID
	l.infrastructure.Title = l.Title
	l.infrastructure.Color = l.Color
	l.infrastructure.ListOptionID = l.ListOptionID
	l.infrastructure.IsHidden = l.IsHidden
}

// reload state from infrastructure.
// It is a mapping function.
func (l *List) reload() error {
	// Get database record and reload when entity is given listID.
	if l.ID != 0 {
		latestList, err := list.FindByID(l.ProjectID, l.ID)
		if err != nil {
			return err
		}
		l.infrastructure = latestList
	}
	l.ID = l.infrastructure.ID
	l.ProjectID = l.infrastructure.ProjectID
	l.UserID = l.infrastructure.UserID
	l.Title = l.infrastructure.Title
	l.Color = l.infrastructure.Color
	l.ListOptionID = l.infrastructure.ListOptionID
	l.IsHidden = l.infrastructure.IsHidden
	return nil
}

// Save call list model save
func (l *List) Save(tx *sql.Tx) error {
	l.reflect()
	return l.infrastructure.Save(tx)
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
	err = l.infrastructure.Update(title, color, listOptionID)
	if err != nil {
		return err
	}
	err = l.reload()
	if err != nil {
		return err
	}
	return nil
}

// Hide call list model hide
func (l *List) Hide() error {
	err := l.infrastructure.Hide()
	if err != nil {
		return err
	}
	return l.reload()
}

// Display call list model display
func (l *List) Display() error {
	err := l.infrastructure.Display()
	if err != nil {
		return err
	}
	return l.reload()
}

// DeleteTasks delete all tasks related a list
func (l *List) DeleteTasks() error {
	return l.infrastructure.DeleteTasks()
}

// Delete delete a list model
func (l *List) Delete() error {
	err := l.infrastructure.Delete()
	if err != nil {
		return err
	}
	l.infrastructure = nil
	return nil
}
