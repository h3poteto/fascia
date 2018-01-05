package inquiry

import (
	"github.com/h3poteto/fascia/server/models/inquiry"
)

type Inquiry struct {
	ID      int64
	Email   string
	Name    string
	Message string
}

func New(id int64, email, name, message string) *Inquiry {
	return &Inquiry{
		id,
		email,
		name,
		message,
	}
}

func (i *Inquiry) Save() (*Inquiry, error) {
	record := inquiry.New(i.ID, i.Email, i.Name, i.Message)
	err := record.Save()
	if err != nil {
		return nil, err
	}
	return New(record.ID, record.Email, record.Name, record.Message), nil
}
