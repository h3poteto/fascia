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

	tables := []string{"tasks", "lists", "projects", "repositories", "reset_passwords", "users", "list_options"}
	// To invalidate foreign key checks for truncate
	// It is valid when database connection pool is not used
	_, err := db.Exec("SET GLOBAL foreign_key_checks = 0;")
	if err != nil {
		return err
	}
	for _, t := range tables {
		_, err := db.Exec("TRUNCATE TABLE " + t + ";")
		if err != nil {
			return errors.Wrapf(err, "truncate failed: %s", t)
		}
	}
	_, err = db.Exec("SET GLOBAL foreign_key_checks = 1;")
	if err != nil {
		return err
	}

	return nil
}

func listOptions() error {
	db := database.SharedInstance().Connection
	_, err := db.Exec("INSERT INTO list_options (action, created_at) values (?, now()), (?, now())",
		"open",
		"close")
	if err != nil {
		return err
	}
	return nil
}
