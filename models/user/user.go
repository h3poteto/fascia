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
	Password string
	Provider string
	OauthToken string
	database db.DB
}

func NewUser(id int64, email string, provider string, oauth_token string) *UserStruct {
	user := &UserStruct{Id: id, Email: email, Provider: provider, OauthToken: oauth_token}
	user.Initialize()
	return user
}

func (u *UserStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}


func CurrentUser(user_id int64) (*UserStruct, error) {
	user := UserStruct{}
	user.Initialize()

	table := user.database.Init()
	defer table.Close()

	id, email, password, provider, oauth_token, created_at, updated_at := int64(0), "", "", "", "", "", ""
	rows, _ := table.Query("select * from users where id = ?;", user_id)
	for rows.Next() {
		err := rows.Scan(&id, &email, &password, &provider, &oauth_token, &created_at, &updated_at)
		if err != nil {
			panic(err.Error())
		}
	}
	user.Id = id
	user.Email = email
	user.Provider = provider
	user.OauthToken = oauth_token
	if id == 0 {
		return &user, errors.New("cannot find user")
	}
	return &user, nil
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

func Login(userEmail string, userPassword string) (*UserStruct, error) {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	id, email, password, provider, oauth_token, created_at, updated_at := int64(0), "", "", "", "", "", ""
	rows, _ := table.Query("select * from users where email = ?;", userEmail)
	for rows.Next() {
		err := rows.Scan(&id, &email, &password, &provider, &oauth_token, &created_at, &updated_at)
		if err != nil {
			panic(err.Error())
		}
	}
	user := NewUser(id, email, provider, oauth_token)
	bytePassword := []byte(userPassword)
	err := bcrypt.CompareHashAndPassword([]byte(password), bytePassword)
	if err != nil {
		fmt.Printf("cannot login: %v\n", userEmail)
		return &UserStruct{}, errors.New("cannot login")
	}
	fmt.Printf("login success\n")
	return user, nil
}

func FindOrCreateGithub(token string) (*UserStruct, error) {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	id, email, password, provider, oauth_token, created_at, updated_at := int64(0), "", "", "", "", "", ""
	rows, _ := table.Query("select * from users where oauth_token = ?;", token)
	for rows.Next() {
		err := rows.Scan(&id, &email, &password, &provider, &oauth_token, &created_at, &updated_at)
		if err != nil {
			panic(err.Error())
		}
	}
	user := NewUser(id, email, provider, oauth_token)
	if id == 0 {
		// email, password更新
		user.Email = "dummy@example.com"
		user.Password = "dummy"
		user.Provider = "github"
		user.OauthToken = token
		user.Save()
	}
	return user, nil

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

func (u *UserStruct) Save() bool {
	table := u.database.Init()
	defer table.Close()

	result, err := table.Exec("insert into users (email, password, provider, oauth_token, created_at) values (?, ?, ?, ?, now());", u.Email, u.Password, u.Provider, u.OauthToken)
	if err != nil {
		return false
	}
	u.Id, _ = result.LastInsertId()
	fmt.Printf("user saved: %v\n", u.Id)
	return true
}
