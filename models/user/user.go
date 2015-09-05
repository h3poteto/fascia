package user

import(
	"../db"
	"time"
	"fmt"
	"golang.org/x/crypto/bcrypt"
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
	u.Initialize()
	table := u.database.Init()
	defer table.Close()

	bytePassword := []byte(password)
	cost := 10
	hashPassword, _ := bcrypt.GenerateFromPassword(bytePassword, cost)
	err := bcrypt.CompareHashAndPassword(hashPassword, bytePassword)
	if err != nil {
		return false
	}
	_, err = table.Exec("insert into users (email, password, created_at) values (?, ?, ?)", email, hashPassword, time.Now())
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
		return false
	}

	return true
}
