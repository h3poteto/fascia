package inquiry

import (
	"github.com/h3poteto/fascia/server/infrastructures/inquiry"
)

// Inquiry is a entity for inquiry.
type Inquiry struct {
	ID             int64
	Email          string
	Name           string
	Message        string
	infrastructure *inquiry.Inquiry
}

// New returns a inquiry struct.
func New(id int64, email, name, message string) *Inquiry {
	infrastructure := inquiry.New(id, email, name, message)
	i := &Inquiry{
		infrastructure: infrastructure,
	}
	i.reload()
	return i
}

func (i *Inquiry) reflect() {
	i.infrastructure.ID = i.ID
	i.infrastructure.Email = i.Email
	i.infrastructure.Name = i.Name
	i.infrastructure.Message = i.Message
}

func (i *Inquiry) reload() error {
	i.ID = i.infrastructure.ID
	i.Email = i.infrastructure.Email
	i.Name = i.infrastructure.Name
	i.Message = i.infrastructure.Message
	return nil
}

// Save a inquiry entity, and returns the latest inquiry object.
func (i *Inquiry) Save() error {
	i.reflect()
	err := i.infrastructure.Save()
	if err != nil {
		return err
	}
	return i.reload()
}
