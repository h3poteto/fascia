package inquiry

import (
	"github.com/h3poteto/fascia/server/models/db"

	"database/sql"

	"github.com/pkg/errors"
)

type Inquiry struct {
	ID       int64
	Email    string
	Name     string
	Message  string
	database *sql.DB
}

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
	i.database = db.SharedInstance().Connection
}

func (i *Inquiry) Save() error {
	result, err := i.database.Exec("insert into inquiries (email, name, message, created_at) values (?, ?, ?, now());", i.Email, i.Name, i.Message)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	i.ID, _ = result.LastInsertId()
	return nil
}
