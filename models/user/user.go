package user

import(
	"../db"
	"time"
	"fmt"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type User interface {
	Registration(string, string) bool
	Login(string, string) (UserStruct, error)
}

type UserStruct struct {
	Id int
	Email string
	database db.DB
}

func NewUser() *UserStruct {
	user := &UserStruct{}
	user.Initialize()
	return user
}

func (u *UserStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

func (u *UserStruct) Registration(email string, password string) bool {
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

func (u *UserStruct) Login(userEmail string, userPassword string) (UserStruct, error) {
	table := u.database.Init()
	defer table.Close()

	id, email, password, created_at, updated_at := 0, "", "", "", ""
	rows, _ := table.Query("select * from users where email = ?;", userEmail)
	for rows.Next() {
		err := rows.Scan(&id, &email, &password, &created_at, &updated_at)
		if err != nil {
			panic(err.Error())
		}
	}
	user := UserStruct{Id: id, Email: email}
	bytePassword := []byte(userPassword)
	err := bcrypt.CompareHashAndPassword([]byte(password), bytePassword)
	if err != nil {
		fmt.Printf("cannot login: %v\n", userEmail)
		return UserStruct{}, errors.New("cannot login")
	}
	fmt.Printf("login success\n")
	return user, nil
}
