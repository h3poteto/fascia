package list_option

import (
	"github.com/h3poteto/fascia/server/infrastructures/list_option"
)

// ListOptionAll list up all options
func ListOptionAll() ([]*ListOption, error) {
	options, err := list_option.All()
	if err != nil {
		return nil, err
	}
	var slice []*ListOption
	for _, option := range options {
		o := &ListOption{
			infrastructure: option,
		}
		if err := o.reload(); err != nil {
			return nil, err
		}
		slice = append(slice, o)
	}
	return slice, nil
}

// FindByID returns a list option
func FindByID(id int64) (*ListOption, error) {
	o := &ListOption{
		ID: id,
	}
	err := o.reload()
	if err != nil {
		return nil, err
	}
	return o, nil
}

// FindByAction returns a list option
func FindByAction(action string) (*ListOption, error) {
	option, err := list_option.FindByAction(action)
	if err != nil {
		return nil, err
	}
	o := &ListOption{
		infrastructure: option,
	}
	if err := o.reload(); err != nil {
		return nil, err
	}
	return o, nil
}

// IsCloseAction return whether it is close option
func (l *ListOption) IsCloseAction() bool {
	if l.Action == "close" {
		return true
	}
	return false
}
