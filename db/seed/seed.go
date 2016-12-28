package seed

import (
	"github.com/h3poteto/fascia/server/models/db"
)

// Seeds insert all seed data
func Seeds() error {
	return listOptions()
}

// TruncateAll remove all records in database
func TruncateAll() error {
	database := db.SharedInstance().Connection

	tables := []string{"list_options", "tasks", "lists", "projects", "repositories", "reset_passwords", "users"}
	for _, t := range tables {
		_, err := database.Exec("TRUNCATE TABLE " + t + ";")
		if err != nil {
			return err
		}
	}

	return nil
}

func listOptions() error {
	database := db.SharedInstance().Connection

	_, err := database.Exec("TRUNCATE TABLE list_options;")
	if err != nil {
		return err
	}
	_, err = database.Exec("INSERT INTO list_options (action, created_at) values (?, now()), (?, now())",
		"open",
		"close")
	if err != nil {
		return err
	}
	return nil
}
