package user

import (
	"../../modules/logging"
	"../db"
	"../project"
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"errors"
	"github.com/google/go-github/github"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"strconv"
	"time"
)

type User interface {
	Projects() []*project.ProjectStruct
}

type UserStruct struct {
	Id         int64
	Email      string
	Password   string
	Provider   sql.NullString
	OauthToken sql.NullString
	Uuid       sql.NullInt64
	UserName   sql.NullString
	Avatar     sql.NullString
	database   db.DB
}

func randomString() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, 36)
}

func HashPassword(password string) ([]byte, error) {
	bytePassword := []byte(password)
	cost := 10
	hashPassword, _ := bcrypt.GenerateFromPassword(bytePassword, cost)
	err := bcrypt.CompareHashAndPassword(hashPassword, bytePassword)
	if err != nil {
		logging.SharedInstance().MethodInfo("user", "hashPassword").Error("hash password error")
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
			logging.SharedInstance().MethodInfo("User", "CurrentUser").Panic(err)
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
		logging.SharedInstance().MethodInfo("user", "CurrentUser").Errorf("cannot find user: %v", userID)
		return &user, errors.New("cannot find user")
	}
	return &user, nil
}

func Registration(email string, password string) (int64, error) {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	hashPassword, err := HashPassword(password)
	if err != nil {
		return 0, err
	}
	result, err := table.Exec("insert into users (email, password, created_at) values (?, ?, ?)", email, hashPassword, time.Now())
	if err != nil {
		logging.SharedInstance().MethodInfo("user", "Registration").Errorf("mysql error: %v", err)
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
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
			logging.SharedInstance().MethodInfo("User", "Login").Panic(err)
		}
	}

	user := NewUser(id, email, provider, oauthToken, uuid, userName, avatarURL)
	bytePassword := []byte(userPassword)
	err := bcrypt.CompareHashAndPassword([]byte(password), bytePassword)
	if err != nil {
		logging.SharedInstance().MethodInfo("user", "Login").Errorf("cannot login: %v", userEmail)
		return nil, errors.New("cannot login")
	}
	logging.SharedInstance().MethodInfo("user", "Login").Info("login success")
	return user, nil
}

func FindUser(id int64) (*UserStruct, error) {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var uuid sql.NullInt64
	var email string
	var provider, oauthToken, userName, avatarURL sql.NullString
	err := table.QueryRow("select email, provider, oauth_token, user_name, uuid, avatar_url from users where id = ?;", id).Scan(&email, &provider, &oauthToken, &userName, &uuid, &avatarURL)
	if err != nil {
		logging.SharedInstance().MethodInfo("User", "FindUser").Infof("cannot find user: %v", err)
		return nil, err
	}
	return NewUser(id, email, provider, oauthToken, uuid, userName, avatarURL), nil
}

func FindByEmail(email string) (*UserStruct, error) {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var id int64
	var uuid sql.NullInt64
	var provider, oauthToken, userName, avatarURL sql.NullString
	err := table.QueryRow("select id, email, provider, oauth_token, user_name, uuid, avatar_url from users where email = ?;", email).Scan(&id, &email, &provider, &oauthToken, &userName, &uuid, &avatarURL)
	if err != nil {
		logging.SharedInstance().MethodInfo("User", "FindByEmail").Infof("cannot find user: %v", err)
		return nil, err
	}
	return NewUser(id, email, provider, oauthToken, uuid, userName, avatarURL), nil
}

// 認証時にもう一度githubアクセスしてidを取ってくるのが無駄なので，できればoauthのcallbakcでidを受け取りたい
func FindOrCreateGithub(token string) (*UserStruct, error) {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	// github認証
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	githubUser, _, _ := client.Users.Get("")

	// TODO: primaryじゃないEmailも保存しておいてログインブロックに使いたい
	emails, _, _ := client.Users.ListEmails(nil)
	var primaryEmail string
	for _, email := range emails {
		if *email.Primary {
			primaryEmail = *email.Email
		}
	}

	var id int64
	var uuid sql.NullInt64
	var email string
	var provider, oauthToken, userName, avatarURL sql.NullString
	rows, _ := table.Query("select id, email, provider, oauth_token, user_name, uuid, avatar_url from users where uuid = ? or email = ?;", *githubUser.ID, primaryEmail)
	for rows.Next() {
		err := rows.Scan(&id, &email, &provider, &oauthToken, &userName, &uuid, &avatarURL)
		if err != nil {
			logging.SharedInstance().MethodInfo("User", "FindOrCreateGithub").Panic(err)
		}
	}
	user := NewUser(id, email, provider, oauthToken, uuid, userName, avatarURL)

	if id == 0 {
		result := user.CreateGithubUser(token, githubUser, primaryEmail)
		if !result {
			logging.SharedInstance().MethodInfo("user", "FindOrCreateGithub").Error("cannot login")
			return user, errors.New("cannot login")
		}
	}

	if !user.OauthToken.Valid || user.OauthToken.String != token {
		result := user.UpdateGithubUserInfo(token, githubUser)
		if !result {
			logging.SharedInstance().MethodInfo("user", "FindOrCreateGithub").Error("cannot update user")
			return user, errors.New("cannot update user")
		}
	}

	return user, nil

}

func (u *UserStruct) Projects() []*project.ProjectStruct {
	table := u.database.Init()
	defer table.Close()

	rows, _ := table.Query("select id, user_id, repository_id, title, description from projects where user_id = ?;", u.Id)
	var slice []*project.ProjectStruct
	for rows.Next() {
		var id, userID int64
		var repositoryID sql.NullInt64
		var title string
		var description string
		err := rows.Scan(&id, &userID, &repositoryID, &title, &description)
		if err != nil {
			logging.SharedInstance().MethodInfo("User", "Projects").Panic(err)
		}
		if id != 0 {
			p := project.NewProject(id, userID, title, description, repositoryID)
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
		logging.SharedInstance().MethodInfo("user", "Save").Errorf("user save error: %v", err)
		return false
	}
	u.Id, _ = result.LastInsertId()
	logging.SharedInstance().MethodInfo("user", "Save").Infof("user saved: %v", u.Id)
	return true
}

func (u *UserStruct) Update() bool {
	table := u.database.Init()
	defer table.Close()

	_, err := table.Exec("update users set provider = ?, oauth_token = ?, uuid = ?, user_name = ?, avatar_url = ? where email = ?;", u.Provider, u.OauthToken, u.Uuid, u.UserName, u.Avatar, u.Email)
	if err != nil {
		logging.SharedInstance().MethodInfo("user", "Update").Errorf("user update error: %v", err)
		return false
	}
	return true
}

func (u *UserStruct) CreateGithubUser(token string, githubUser *github.User, primaryEmail string) bool {
	u.Email = primaryEmail
	bytePassword, _ := HashPassword(randomString())
	u.Password = string(bytePassword)
	u.Provider = sql.NullString{String: "github", Valid: true}
	u.OauthToken = sql.NullString{String: token, Valid: true}

	u.UserName = sql.NullString{String: *githubUser.Login, Valid: true}
	u.Uuid = sql.NullInt64{Int64: int64(*githubUser.ID), Valid: true}
	u.Avatar = sql.NullString{String: *githubUser.AvatarURL, Valid: true}
	u.Save()
	return true
}

func (u *UserStruct) UpdateGithubUserInfo(token string, githubUser *github.User) bool {
	u.Provider = sql.NullString{String: "github", Valid: true}
	u.OauthToken = sql.NullString{String: token, Valid: true}
	u.UserName = sql.NullString{String: *githubUser.Login, Valid: true}
	u.Uuid = sql.NullInt64{Int64: int64(*githubUser.ID), Valid: true}
	u.Avatar = sql.NullString{String: *githubUser.AvatarURL, Valid: true}
	u.Update()
	return true
}
