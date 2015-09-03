package db
import (
	"os"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DB interface {
	Init() *sql.DB
}

type Database struct {
}

func (u *Database) Init() *sql.DB {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	db, err := sql.Open("mysql", username + ":" + password + "@/fascia?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	return db
}
