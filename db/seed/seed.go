package seed

import (
	"github.com/h3poteto/fascia/server/models/db"
)

// Seeds insert all seed data
func Seeds() error {
	return listOptions()
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
