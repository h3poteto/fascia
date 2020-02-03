package list

// NewOption returns a list option entity
func NewOption(id int64, action string) *Option {
	return &Option{
		id,
		action,
	}
}

// IsCloseAction return whether it is close option
func (o *Option) IsCloseAction() bool {
	if o == nil {
		return false
	}
	if o.Action == "close" {
		return true
	}
	return false
}
