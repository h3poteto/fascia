package main

import (
	"../../models/db"
)

func main() {
	ListOptions()
}

func ListOptions() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()

	_, err := table.Exec("TRUNCATE TABLE list_options;")
	if err != nil {
		panic(err)
	}
	_, err = table.Exec("INSERT INTO list_options (action, created_at) values (?, now()), (?, now())",
		"open",
		"close")
	if err != nil {
		panic(err)
	}
}
