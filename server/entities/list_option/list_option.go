package list_option

import (
	"github.com/h3poteto/fascia/server/infrastructures/list_option"
)

// ListOption has a list option model object
type ListOption struct {
	ID             int64
	Action         string
	infrastructure *list_option.ListOption
}

// New returns a list option entity
func New(id int64, action string) *ListOption {
	infrastructure := list_option.New(id, action)
	o := &ListOption{
		infrastructure: infrastructure,
	}
	o.reload()
	return o
}

func (o *ListOption) reflect() {
	o.infrastructure.ID = o.ID
	o.infrastructure.Action = o.Action
}

func (o *ListOption) reload() error {
	if o.ID != 0 {
		latestOption, err := list_option.FindByID(o.ID)
		if err != nil {
			return err
		}
		o.infrastructure = latestOption
	}
	o.ID = o.infrastructure.ID
	o.Action = o.infrastructure.Action
	return nil
}
