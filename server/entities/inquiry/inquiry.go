package inquiry

import (
	"github.com/h3poteto/fascia/server/infrastructures/inquiry"
)

// Inquiry is a entity for inquiry.
type Inquiry struct {
	ID      int64
	Email   string
	Name    string
	Message string
}

// New returns a inquiry struct.
func New(id int64, email, name, message string) *Inquiry {
	return &Inquiry{
		id,
		email,
		name,
		message,
	}
}

// Save a inquiry entity, and returns the latest inquiry object.
func (i *Inquiry) Save() (*Inquiry, error) {
	record := inquiry.New(i.ID, i.Email, i.Name, i.Message)
	err := record.Save()
	if err != nil {
		return nil, err
	}
	return New(record.ID, record.Email, record.Name, record.Message), nil
}
