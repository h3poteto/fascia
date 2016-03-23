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
	ID         int64
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
	hashPassword, err := bcrypt.GenerateFromPassword(bytePassword, cost)
	if err != nil {
		logging.SharedInstance().MethodInfo("user", "hashPassword", true).Errorf("generate password error: %v", err)
		return nil, errors.New("generate password error")
	}
	err = bcrypt.CompareHashAndPassword(hashPassword, bytePassword)
	if err != nil {
		logging.SharedInstance().MethodInfo("user", "hashPassword", true).Errorf("hash password error: %v", err)
		return nil, errors.New("hash password error")
	}
	return hashPassword, nil
}

func NewUser(id int64, email string, provider sql.NullString, oauthToken sql.NullString, uuid sql.NullInt64, userName sql.NullString, avatar sql.NullString) *UserStruct {
	user := &UserStruct{ID: id, Email: email, Provider: provider, OauthToken: oauthToken, Uuid: uuid, UserName: userName, Avatar: avatar}
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
	err := table.QueryRow("select id, email, provider, oauth_token, user_name, uuid, avatar_url from users where id = ?;", userID).Scan(&id, &email, &provider, &oauthToken, &userName, &uuid, &avatarURL)
	if err != nil {
		logging.SharedInstance().MethodInfo("User", "CurrentUser").Infof("cannot find user: %v", err)
		return nil, err
	}
	user.ID = id
	user.Email = email
	user.Provider = provider
	user.OauthToken = oauthToken
	user.UserName = userName
	user.Uuid = uuid
	user.Avatar = avatarURL
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
		logging.SharedInstance().MethodInfo("user", "Registration", true).Errorf("mysql error: %v", err)
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
	err := table.QueryRow("select id, email, password, provider, oauth_token, user_name, uuid, avatar_url from users where email = ?;", userEmail).Scan(&id, &email, &password, &provider, &oauthToken, &userName, &uuid, &avatarURL)
	if err != nil {
		logging.SharedInstance().MethodInfo("User", "Login").Infof("cannot find user: %v", err)
		return nil, err
	}

	user := NewUser(id, email, provider, oauthToken, uuid, userName, avatarURL)
	bytePassword := []byte(userPassword)
	err = bcrypt.CompareHashAndPassword([]byte(password), bytePassword)
	if err != nil {
		logging.SharedInstance().MethodInfo("user", "Login").Debugf("cannot login: %v", userEmail)
		return nil, errors.New("cannot login")
	}
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
	githubUser, _, err := client.Users.Get("")
	if err != nil {
		logging.SharedInstance().MethodInfo("User", "FindOrCreateGithub", true).Errorf("github user find error: %v", err)
		return nil, errors.New("github user find error")
	}

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
	rows, err := table.Query("select id, email, provider, oauth_token, user_name, uuid, avatar_url from users where uuid = ? or email = ?;", *githubUser.ID, primaryEmail)
	if err != nil {
		logging.SharedInstance().MethodInfo("User", "FindOrCreateGithub", true).Panic(err)
	}
	for rows.Next() {
		err := rows.Scan(&id, &email, &provider, &oauthToken, &userName, &uuid, &avatarURL)
		if err != nil {
			logging.SharedInstance().MethodInfo("User", "FindOrCreateGithub", true).Panic(err)
		}
	}
	user := NewUser(id, email, provider, oauthToken, uuid, userName, avatarURL)

	if id == 0 {
		if err := user.CreateGithubUser(token, githubUser, primaryEmail); err != nil {
			logging.SharedInstance().MethodInfo("user", "FindOrCreateGithub", true).Error("cannot login to github")
			return user, errors.New("cannot login to github")
		}
	}

	if !user.OauthToken.Valid || user.OauthToken.String != token {
		if err := user.UpdateGithubUserInfo(token, githubUser); err != nil {
			logging.SharedInstance().MethodInfo("user", "FindOrCreateGithub", true).Error("cannot update user")
			return user, errors.New("cannot update user")
		}
	}

	return user, nil

}

func (u *UserStruct) Projects() []*project.ProjectStruct {
	table := u.database.Init()
	defer table.Close()

	var slice []*project.ProjectStruct
	rows, err := table.Query("select id, user_id, repository_id, title, description, show_issues, show_pull_requests from projects where user_id = ?;", u.ID)
	if err != nil {
		logging.SharedInstance().MethodInfo("User", "Projects", true).Panic(err)
	}
	for rows.Next() {
		var id, userID int64
		var repositoryID sql.NullInt64
		var title string
		var description string
		var showIssues, showPullRequests bool
		err := rows.Scan(&id, &userID, &repositoryID, &title, &description, &showIssues, &showPullRequests)
		if err != nil {
			logging.SharedInstance().MethodInfo("User", "Projects", true).Panic(err)
		}
		if id != 0 {
			p := project.NewProject(id, userID, title, description, repositoryID, showIssues, showPullRequests)
			slice = append(slice, p)
		}
	}
	return slice
}

func (u *UserStruct) Save() error {
	table := u.database.Init()
	defer table.Close()

	result, err := table.Exec("insert into users (email, password, provider, oauth_token, uuid, user_name, avatar_url, created_at) values (?, ?, ?, ?, ?, ?, ?, now());", u.Email, u.Password, u.Provider, u.OauthToken, u.Uuid, u.UserName, u.Avatar)
	if err != nil {
		logging.SharedInstance().MethodInfo("user", "Save", true).Errorf("failed to create user: %v", err)
		return err
	}
	u.ID, _ = result.LastInsertId()
	logging.SharedInstance().MethodInfo("user", "Save").Infof("user saved: %v", u.ID)
	return nil
}

func (u *UserStruct) Update() error {
	table := u.database.Init()
	defer table.Close()

	_, err := table.Exec("update users set provider = ?, oauth_token = ?, uuid = ?, user_name = ?, avatar_url = ? where email = ?;", u.Provider, u.OauthToken, u.Uuid, u.UserName, u.Avatar, u.Email)
	if err != nil {
		logging.SharedInstance().MethodInfo("User", "Update", true).Panic(err)
		return err
	}
	return nil
}

func (u *UserStruct) CreateGithubUser(token string, githubUser *github.User, primaryEmail string) error {
	u.Email = primaryEmail
	bytePassword, err := HashPassword(randomString())
	if err != nil {
		logging.SharedInstance().MethodInfo("User", "CreateGithubUser").Error(err)
		return err
	}
	u.Password = string(bytePassword)
	u.Provider = sql.NullString{String: "github", Valid: true}
	u.OauthToken = sql.NullString{String: token, Valid: true}

	u.UserName = sql.NullString{String: *githubUser.Login, Valid: true}
	u.Uuid = sql.NullInt64{Int64: int64(*githubUser.ID), Valid: true}
	u.Avatar = sql.NullString{String: *githubUser.AvatarURL, Valid: true}
	if err := u.Save(); err != nil {
		return err
	}
	return nil
}

func (u *UserStruct) UpdateGithubUserInfo(token string, githubUser *github.User) error {
	u.Provider = sql.NullString{String: "github", Valid: true}
	u.OauthToken = sql.NullString{String: token, Valid: true}
	u.UserName = sql.NullString{String: *githubUser.Login, Valid: true}
	u.Uuid = sql.NullInt64{Int64: int64(*githubUser.ID), Valid: true}
	u.Avatar = sql.NullString{String: *githubUser.AvatarURL, Valid: true}
	if err := u.Update(); err != nil {
		return err
	}
	return nil
}
