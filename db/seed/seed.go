package main

import (
	"../../models/db"
)

func main() {
	list_options()
}

func list_options() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()

	_, err := table.Exec("TRUNCATE TABLE list_options;")
	if err != nil {
		panic(err.Error())
	}
	_, err = table.Exec("INSERT INTO list_options (action, created_at) values (?, now()), (?, now())",
		"reopen",
		"close")
	if err != nil {
		panic(err.Error())
	}
}
