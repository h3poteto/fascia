package list_option

import (
	"../db"
	"database/sql"
)

type ListOiption interface {
}

type ListOptionStruct struct {
	Id       int64
	Action   string
	database db.DB
}
