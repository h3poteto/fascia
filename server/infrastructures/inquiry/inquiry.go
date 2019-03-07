package inquiry

import (
	"database/sql"

	"github.com/pkg/errors"
)

// Inquiry is a record object for inquiry.
type Inquiry struct {
	db *sql.DB
}

// New returns a inquiry struct.
func New(db *sql.DB) *Inquiry {
	inquiry := &Inquiry{
		db,
	}
	return inquiry
}

// Create a inquiry object in database.
func (i *Inquiry) Create(email, name, message string) (int64, error) {
	result, err := i.db.Exec("insert into inquiries (email, name, message, created_at) values (?, ?, ?, now());", email, name, message)
	if err != nil {
		return 0, errors.Wrap(err, "inquiry repository")
	}
	id, _ := result.LastInsertId()
	return id, nil
}
