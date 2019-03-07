package list

import "github.com/pkg/errors"

// Option has a list option model object
type Option struct {
	ID             int64
	Action         string
	infrastructure Repository
}

// NewOption returns a list option entity
func NewOption(id int64, action string, infrastructure Repository) *Option {
	return &Option{
		id,
		action,
		infrastructure,
	}
}

// IsCloseAction return whether it is close option
func (o *Option) IsCloseAction() bool {
	if o.Action == "close" {
		return true
	}
	return false
}

// AllOption list up all options
func AllOption(infrastructure Repository) ([]*Option, error) {
	options, err := infrastructure.AllOption()
	if err != nil {
		return nil, err
	}

	result := []*Option{}
	for _, option := range options {
		id, ok := option["id"].(int64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		action, ok := option["action"].(string)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		l := NewOption(id, action, infrastructure)
		result = append(result, l)
	}
	return result, nil
}

// FindOptionByID returns a list option
func FindOptionByID(id int64, infrastructure Repository) (*Option, error) {
	if id == 0 {
		return nil, errors.New("Please set option id")
	}
	id, action, err := infrastructure.FindOptionByID(id)
	if err != nil {
		return nil, err
	}
	return NewOption(id, action, infrastructure), nil
}

// FindOptionByAction returns a list option
func FindOptionByAction(action string, infrastructure Repository) (*Option, error) {
	id, action, err := infrastructure.FindOptionByAction(action)
	if err != nil {
		return nil, err
	}
	return NewOption(id, action, infrastructure), nil
}
