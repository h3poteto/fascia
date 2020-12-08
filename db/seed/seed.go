package seed

import (
	"github.com/h3poteto/fascia/lib/modules/database"
	"github.com/pkg/errors"
)

// Seeds insert all seed data
func Seeds() error {
	return listOptions()
}

// TruncateAll remove all records in database
func TruncateAll() error {
	db := database.SharedInstance().Connection

	tables := []string{"tasks", "lists", "projects", "repositories", "users", "list_options"}
	for _, t := range tables {
		_, err := db.Exec("DELETE FROM " + t + ";")
		if err != nil {
			return errors.Wrapf(err, "truncate failed: %s", t)
		}
	}

	return nil
}

func listOptions() error {
	db := database.SharedInstance().Connection
	_, err := db.Exec("INSERT INTO list_options (action) values ($1), ($2)",
		"open",
		"close")
	if err != nil {
		return err
	}
	return nil
}
