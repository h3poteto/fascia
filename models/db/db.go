package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type DB interface {
	Init() *sql.DB
}

type Database struct {
}

func (u *Database) Init() *sql.DB {
	env := os.Getenv("GOJIENV")
	buf, err := ioutil.ReadFile("db/dbconf.yml")
	if err != nil {
		panic(err)
	}
	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		panic(err)
	}
	username := m[env].(map[interface{}]interface{})["user"].(string)
	password := m[env].(map[interface{}]interface{})["password"].(string)
	database := m[env].(map[interface{}]interface{})["name"].(string)
	username = os.ExpandEnv(username)
	password = os.ExpandEnv(password)
	database = os.ExpandEnv(database)
	db, err := sql.Open("mysql", username+":"+password+"@/"+database+"?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	return db
}
