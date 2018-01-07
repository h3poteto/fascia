package inquiry

import (
	"github.com/h3poteto/fascia/lib/modules/database"

	"database/sql"

	"github.com/pkg/errors"
)

// Inquiry is a record object for inquiry.
type Inquiry struct {
	ID      int64
	Email   string
	Name    string
	Message string
	db      *sql.DB
}

// New returns a inquiry struct.
func New(id int64, email, name, message string) *Inquiry {
	inquiry := &Inquiry{
		ID:      id,
		Email:   email,
		Name:    name,
		Message: message,
	}
	inquiry.initialize()
	return inquiry
}

func (i *Inquiry) initialize() {
	i.db = database.SharedInstance().Connection
}

// Save a inquiry object in database.
func (i *Inquiry) Save() error {
	result, err := i.db.Exec("insert into inquiries (email, name, message, created_at) values (?, ?, ?, now());", i.Email, i.Name, i.Message)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	i.ID, _ = result.LastInsertId()
	return nil
}