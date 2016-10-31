package main

import (
	"../../models/db"
)

func main() {
	ListOptions()
}

func ListOptions() {
	database := db.SharedInstance().Connection

	_, err := database.Exec("TRUNCATE TABLE list_options;")
	if err != nil {
		panic(err)
	}
	_, err = database.Exec("INSERT INTO list_options (action, created_at) values (?, now()), (?, now())",
		"open",
		"close")
	if err != nil {
		panic(err)
	}
}
