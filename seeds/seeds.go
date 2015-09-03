package main

import (
	"log"
	_ "github.com/go-sql-driver/mysql"

	"../models/db"
)

func main() {
	create_users()
}


func create_users() {
	objectDB := &db.Database{}
	var database db.DB = objectDB
	mydb := database.Init()
	defer mydb.Close()

	_, err := mydb.Query("select users.id from users;")
	if err != nil {
		_, err = mydb.Exec("CREATE TABLE users (id int(11) NOT NULL AUTO_INCREMENT, email varchar(255) DEFAULT NULL, password varchar(255)  DEFAULT NULL, created_at datetime DEFAULT NULL, updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (id)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;")
		if err != nil {
			log.Fatalf("mysql error: %v ", err)
		}
	}
}
