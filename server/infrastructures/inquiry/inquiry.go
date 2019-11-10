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
	var id int64
	err := i.db.QueryRow("INSERT INTO inquiries (email, name, message) VALUES ($1, $2, $3) RETURNING id;", email, name, message).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "inquiry repository")
	}
	return id, nil
}
