package user

import(
	"../db"
	"time"
	"fmt"
	"errors"
	"database/sql"
	"encoding/binary"
	"crypto/rand"
	"strconv"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
	"../project"
)

type User interface {
	Projects() []*project.ProjectStruct
}

type UserStruct struct {
	Id int64
	Email string
	Password string
	Provider sql.NullString
	OauthToken sql.NullString
	Uuid sql.NullInt64
	UserName sql.NullString
	Avatar sql.NullString
	database db.DB
}

func randomString() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, 36)
}

func hashPassword(password string) ([]byte, error) {
	bytePassword := []byte(password)
	cost := 10
	hashPassword, _ := bcrypt.GenerateFromPassword(bytePassword, cost)
	err := bcrypt.CompareHashAndPassword(hashPassword, bytePassword)
	if err != nil {
		return nil, errors.New("hash password error")
	}
	return hashPassword, nil
}


func NewUser(id int64, email string, provider sql.NullString, oauthToken sql.NullString, uuid sql.NullInt64, userName sql.NullString, avatar sql.NullString) *UserStruct {
	user := &UserStruct{Id: id, Email: email, Provider: provider, OauthToken: oauthToken, Uuid: uuid, UserName: userName, Avatar: avatar}
	user.Initialize()
	return user
}

func (u *UserStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}


func CurrentUser(userID int64) (*UserStruct, error) {
	user := UserStruct{}
	user.Initialize()

	table := user.database.Init()
	defer table.Close()

	var id int64
	var uuid sql.NullInt64
	var email string
	var provider, oauthToken, userName, avatarURL sql.NullString
	rows, _ := table.Query("select id, email, provider, oauth_token, user_name, uuid, avatar_url from users where id = ?;", userID)
	for rows.Next() {
		err := rows.Scan(&id, &email, &provider, &oauthToken, &userName, &uuid, &avatarURL)
		if err != nil {
			panic(err.Error())
		}
	}
	user.Id = id
	user.Email = email
	user.Provider = provider
	user.OauthToken = oauthToken
	user.UserName = userName
	user.Uuid = uuid
	user.Avatar = avatarURL
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

	hashPassword, err := hashPassword(password)
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

	var id int64
	var uuid sql.NullInt64
	var email, password string
	var provider, oauthToken, userName, avatarURL sql.NullString
	rows, _ := table.Query("select id, email, password, provider, oauth_token, user_name, uuid, avatar_url from users where email = ?;", userEmail)
	for rows.Next() {
		err := rows.Scan(&id, &email, &password, &provider, &oauthToken, &userName, &uuid, &avatarURL)
		if err != nil {
			panic(err.Error())
		}
	}

	user := NewUser(id, email, provider, oauthToken, uuid, userName, avatarURL)
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

	var id int64
	var uuid sql.NullInt64
	var email string
	var provider, oauthToken, userName, avatarURL sql.NullString
	rows, _ := table.Query("select id, email, provider, oauth_token, user_name, uuid, avatar_url from users where oauth_token = ?;", token)
	for rows.Next() {
		err := rows.Scan(&id, &email, &provider, &oauthToken, &userName, &uuid, &avatarURL)
		if err != nil {
			panic(err.Error())
		}
	}
	user := NewUser(id, email, provider, oauthToken, uuid, userName, avatarURL)

 	if id == 0 {
		user.CreateGithubUser(token)
	}
	return user, nil

}

func (u *UserStruct) Projects() []*project.ProjectStruct {
	table := u.database.Init()
	defer table.Close()

	rows, _ := table.Query("select id, user_id, title from projects where user_id = ?;", u.Id)
	var slice []*project.ProjectStruct
	for rows.Next() {
		var id int64
		var userID sql.NullInt64
		var title string
		err := rows.Scan(&id, &userID, &title)
		if err != nil {
			panic(err.Error())
		}
		if userID.Valid {
			p := project.NewProject(id, userID.Int64, title)
			slice = append(slice, p)
		}
	}
	return slice
}

func (u *UserStruct) Save() bool {
	table := u.database.Init()
	defer table.Close()

	result, err := table.Exec("insert into users (email, password, provider, oauth_token, uuid, user_name, avatar_url, created_at) values (?, ?, ?, ?, ?, ?, ?, now());", u.Email, u.Password, u.Provider, u.OauthToken, u.Uuid, u.UserName, u.Avatar)
	if err != nil {
		fmt.Printf("login error: %+v\n", err)
		return false
	}
	u.Id, _ = result.LastInsertId()
	fmt.Printf("user saved: %v\n", u.Id)
	return true
}

func (u *UserStruct) CreateGithubUser(token string) bool {
	// email, password更新
	// TODO: ここuniqとってるのでもっと慎重にアドレス決定しないとやばい
	u.Email = randomString() + "@fascia.io"
	bytePassword, _ := hashPassword(randomString())
	u.Password = string(bytePassword)
	u.Provider = sql.NullString{String: "github", Valid: true}
	u.OauthToken = sql.NullString{String: token, Valid: true}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	user, _, _ := client.Users.Get("")
	u.UserName = sql.NullString{String: *user.Login, Valid: true}
	u.Uuid = sql.NullInt64{Int64: int64(*user.ID), Valid: true}
	u.Avatar = sql.NullString{String: *user.AvatarURL, Valid: true}
	u.Save()
	return true
}
