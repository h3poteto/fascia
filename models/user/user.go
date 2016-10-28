package user

import (
	"../../modules/logging"
	"../db"
	"../project"

	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
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
	database   *sql.DB
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
		return nil, errors.Wrap(err, "cannot generate password")
	}
	err = bcrypt.CompareHashAndPassword(hashPassword, bytePassword)
	if err != nil {
		return nil, errors.Wrap(err, "did not match password")
	}
	return hashPassword, nil
}

func NewUser(id int64, email string, provider sql.NullString, oauthToken sql.NullString, uuid sql.NullInt64, userName sql.NullString, avatar sql.NullString) *UserStruct {
	user := &UserStruct{ID: id, Email: email, Provider: provider, OauthToken: oauthToken, Uuid: uuid, UserName: userName, Avatar: avatar}
	user.Initialize()
	return user
}

func (u *UserStruct) Initialize() {
	u.database = db.SharedInstance().Connection
}

func CurrentUser(userID int64) (*UserStruct, error) {
	user := UserStruct{}
	user.Initialize()

	var id int64
	var uuid sql.NullInt64
	var email string
	var provider, oauthToken, userName, avatarURL sql.NullString
	err := user.database.QueryRow("select id, email, provider, oauth_token, user_name, uuid, avatar_url from users where id = ?;", userID).Scan(&id, &email, &provider, &oauthToken, &userName, &uuid, &avatarURL)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
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

// Registration is create new user through validation
func Registration(email string, password string, passwordConfirm string) (int64, error) {
	hashPassword, err := HashPassword(password)
	if err != nil {
		return 0, err
	}

	user := NewUser(0, email, sql.NullString{}, sql.NullString{}, sql.NullInt64{}, sql.NullString{}, sql.NullString{})
	user.Password = string(hashPassword)
	err = user.Save()

	return user.ID, err
}

func Login(userEmail string, userPassword string) (*UserStruct, error) {
	database := db.SharedInstance().Connection
	var id int64
	var uuid sql.NullInt64
	var email, password string
	var provider, oauthToken, userName, avatarURL sql.NullString
	err := database.QueryRow("select id, email, password, provider, oauth_token, user_name, uuid, avatar_url from users where email = ?;", userEmail).Scan(&id, &email, &password, &provider, &oauthToken, &userName, &uuid, &avatarURL)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}

	user := NewUser(id, email, provider, oauthToken, uuid, userName, avatarURL)
	bytePassword := []byte(userPassword)
	err = bcrypt.CompareHashAndPassword([]byte(password), bytePassword)
	if err != nil {
		return nil, errors.Wrap(err, "password did not match")
	}
	return user, nil
}

func FindUser(id int64) (*UserStruct, error) {
	database := db.SharedInstance().Connection

	var uuid sql.NullInt64
	var email string
	var provider, oauthToken, userName, avatarURL sql.NullString
	err := database.QueryRow("select email, provider, oauth_token, user_name, uuid, avatar_url from users where id = ?;", id).Scan(&email, &provider, &oauthToken, &userName, &uuid, &avatarURL)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	return NewUser(id, email, provider, oauthToken, uuid, userName, avatarURL), nil
}

func FindByEmail(email string) (*UserStruct, error) {
	database := db.SharedInstance().Connection
	var id int64
	var uuid sql.NullInt64
	var provider, oauthToken, userName, avatarURL sql.NullString
	err := database.QueryRow("select id, email, provider, oauth_token, user_name, uuid, avatar_url from users where email = ?;", email).Scan(&id, &email, &provider, &oauthToken, &userName, &uuid, &avatarURL)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	return NewUser(id, email, provider, oauthToken, uuid, userName, avatarURL), nil
}

// 認証時にもう一度githubアクセスしてidを取ってくるのが無駄なので，できればoauthのcallbakcでidを受け取りたい
func FindOrCreateGithub(token string) (*UserStruct, error) {
	// github認証
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	githubUser, _, err := client.Users.Get("")
	if err != nil {
		return nil, errors.Wrap(err, "github api error")
	}

	// TODO: primaryじゃないEmailも保存しておいてログインブロックに使いたい
	emails, _, _ := client.Users.ListEmails(nil)
	var primaryEmail string
	for _, email := range emails {
		if *email.Primary {
			primaryEmail = *email.Email
		}
	}

	database := db.SharedInstance().Connection
	var id int64
	var uuid sql.NullInt64
	var email string
	var provider, oauthToken, userName, avatarURL sql.NullString
	rows, err := database.Query("select id, email, provider, oauth_token, user_name, uuid, avatar_url from users where uuid = ? or email = ?;", *githubUser.ID, primaryEmail)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	for rows.Next() {
		err := rows.Scan(&id, &email, &provider, &oauthToken, &userName, &uuid, &avatarURL)
		if err != nil {
			return nil, errors.Wrap(err, "sql scan error")
		}
	}
	user := NewUser(id, email, provider, oauthToken, uuid, userName, avatarURL)

	if id == 0 {
		if err := user.CreateGithubUser(token, githubUser, primaryEmail); err != nil {
			return user, err
		}
	}

	if !user.OauthToken.Valid || user.OauthToken.String != token {
		if err := user.UpdateGithubUserInfo(token, githubUser); err != nil {
			return user, err
		}
	}

	return user, nil

}

// Projects list up projects related a user
func (u *UserStruct) Projects() ([]*project.ProjectStruct, error) {
	var slice []*project.ProjectStruct
	rows, err := u.database.Query("select id, user_id, repository_id, title, description, show_issues, show_pull_requests from projects where user_id = ?;", u.ID)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	for rows.Next() {
		var id, userID int64
		var repositoryID sql.NullInt64
		var title string
		var description string
		var showIssues, showPullRequests bool
		err := rows.Scan(&id, &userID, &repositoryID, &title, &description, &showIssues, &showPullRequests)
		if err != nil {
			return nil, errors.Wrap(err, "sql select error")
		}
		if id != 0 {
			p := project.NewProject(id, userID, title, description, repositoryID, showIssues, showPullRequests)
			slice = append(slice, p)
		}
	}
	return slice, nil
}

func (u *UserStruct) Save() error {
	// TODO: この前にvalidationを入れたい
	result, err := u.database.Exec("insert into users (email, password, provider, oauth_token, uuid, user_name, avatar_url, created_at) values (?, ?, ?, ?, ?, ?, ?, now());", u.Email, u.Password, u.Provider, u.OauthToken, u.Uuid, u.UserName, u.Avatar)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	u.ID, _ = result.LastInsertId()
	logging.SharedInstance().MethodInfo("user", "Save").Infof("user saved: %v", u.ID)
	return nil
}

func (u *UserStruct) Update() error {
	_, err := u.database.Exec("update users set provider = ?, oauth_token = ?, uuid = ?, user_name = ?, avatar_url = ? where email = ?;", u.Provider, u.OauthToken, u.Uuid, u.UserName, u.Avatar, u.Email)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	return nil
}

func (u *UserStruct) CreateGithubUser(token string, githubUser *github.User, primaryEmail string) error {
	u.Email = primaryEmail
	bytePassword, err := HashPassword(randomString())
	if err != nil {
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
