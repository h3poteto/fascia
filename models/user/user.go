package user

import(
	"../db"
	"time"
	"fmt"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"../project"
)

type User interface {
	Projects() []*project.ProjectStruct
}

type UserStruct struct {
	Id int64
	Email string
	database db.DB
}

func NewUser(id int64, email string) *UserStruct {
	user := &UserStruct{Id: id, Email: email}
	user.Initialize()
	return user
}

func (u *UserStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}


func CurrentUser(user_id int64) (UserStruct, error) {
	user := UserStruct{}
	user.Initialize()

	table := user.database.Init()
	defer table.Close()

	id, email, password, created_at, updated_at := int64(0), "", "", "", ""
	rows, _ := table.Query("select * from users where id = ?;", user_id)
	for rows.Next() {
		err := rows.Scan(&id, &email, &password, &created_at, &updated_at)
		if err != nil {
			panic(err.Error())
		}
	}
	user.Id = id
	user.Email = email
	if id == 0 {
		return user, errors.New("cannot find user")
	}
	return user, nil
}

func Registration(email string, password string) bool {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
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

func Login(userEmail string, userPassword string) (UserStruct, error) {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	id, email, password, created_at, updated_at := int64(0), "", "", "", ""
	rows, _ := table.Query("select * from users where email = ?;", userEmail)
	for rows.Next() {
		err := rows.Scan(&id, &email, &password, &created_at, &updated_at)
		if err != nil {
			panic(err.Error())
		}
	}
	user := NewUser(id, email)
	bytePassword := []byte(userPassword)
	err := bcrypt.CompareHashAndPassword([]byte(password), bytePassword)
	if err != nil {
		fmt.Printf("cannot login: %v\n", userEmail)
		return UserStruct{}, errors.New("cannot login")
	}
	fmt.Printf("login success\n")
	return *user, nil
}

func (u *UserStruct) Projects() []*project.ProjectStruct {
	table := u.database.Init()
	defer table.Close()

	rows, _ := table.Query("select * from projects where user_id = ?;", u.Id)
	var slice []*project.ProjectStruct
	for rows.Next() {
		id, user_id, title, created_at, updated_at := int64(0), int64(0), "", "", ""
		err := rows.Scan(&id, &user_id, &title, &created_at, &updated_at)
		if err != nil {
			panic(err.Error())
		}
		p := project.NewProject(id, user_id, title)
		slice = append(slice, p)
	}
	return slice
}
