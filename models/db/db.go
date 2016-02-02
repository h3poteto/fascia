package db

import (
	"../../modules/logging"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type DB interface {
	Init() *sql.DB
}

type Database struct {
}

func (u *Database) Init() *sql.DB {
	env := os.Getenv("GOJIENV")
	root := os.Getenv("GOJIROOT")
	path := filepath.Join(root, "db/dbconf.yml")
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		logging.SharedInstance().MethodInfo("DB", "Init", true).Panic(err)
	}
	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		logging.SharedInstance().MethodInfo("DB", "Init", true).Panic(err)
	}
	username := m[env].(map[interface{}]interface{})["user"].(string)
	password := m[env].(map[interface{}]interface{})["password"].(string)
	database := m[env].(map[interface{}]interface{})["name"].(string)
	username = os.ExpandEnv(username)
	password = os.ExpandEnv(password)
	database = os.ExpandEnv(database)
	db, err := sql.Open("mysql", username+":"+password+"@/"+database+"?charset=utf8")
	if err != nil {
		logging.SharedInstance().MethodInfo("DB", "Init", true).Panic(err)
	}
	return db
}
