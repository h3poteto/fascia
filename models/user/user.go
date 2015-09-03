package user

import(
	"../db"
	"time"
	"fmt"
)

type User interface {
	Initialize()
	Registration(string, string) bool
}

type UserStruct struct {
	database db.DB
}

func (u *UserStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

func (u *UserStruct) Registration(email string, password string) bool {
	table := u.database.Init()
	defer table.Close()

	_, err := table.Exec("insert into users (email, password, created_at, updated_at) values (?, ?, ?, ?)", email, password, time.Now(), time.Now())
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
		return false
	}

	return true
}
