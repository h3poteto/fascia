package database

import (
	"database/sql"
	"io/ioutil"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

type Database struct {
	Connection *sql.DB
}

var sharedInstance = New()

func New() *Database {
	env := os.Getenv("APPENV")
	path := "db/dbconf.yml"
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		panic(err)
	}
	open := m[env].(map[interface{}]interface{})["open"].(string)
	pool := m[env].(map[interface{}]interface{})["pool"].(int)
	open = os.ExpandEnv(open)

	db, err := sql.Open("postgres", open)
	if err != nil {
		panic(err)
	}

	// MaxIdle: dbへのアクセスがないときにも保持しておくconnection poolの上限
	// MaxOpen: idle + activeなconnection poolの上限数
	db.SetMaxIdleConns(pool)
	db.SetMaxOpenConns(pool)

	return &Database{
		Connection: db,
	}
}

func SharedInstance() *Database {
	return sharedInstance
}

func (d *Database) Close() error {
	return d.Connection.Close()
}
