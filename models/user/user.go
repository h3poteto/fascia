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
	"regexp"
	"strconv"
	"unicode/utf8"
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

// TODO: ここのバリデーションはむしろformバリデーションなので，モデル層で持ちたくない
// validator層を作るべき
func emailValidation(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return re.MatchString(email)
}

// HashPassword mask password
func HashPassword(password string) []byte {
	bytePassword := []byte(password)
	cost := 10
	hashPassword, err := bcrypt.GenerateFromPassword(bytePassword, cost)
	if err != nil {
		panic(err)
	}
	err = bcrypt.CompareHashAndPassword(hashPassword, bytePassword)
	if err != nil {
		panic(err)
	}
	return hashPassword
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

// Validation is check email and password for user model
func Validation(email string, password string, passwordConfirm string) bool {
	if !emailValidation(email) {
		return false
	}
	if password != passwordConfirm {
		return false
	}
	if utf8.RuneCountInString(password) < 8 {
		return false
	}
	return true
}

// Registration is create new user through validation
func Registration(email string, password string, passwordConfirm string) (int64, error) {
	if !Validation(email, password, passwordConfirm) {
		return 0, errors.New("validation failed")
	}

	user := NewUser(0, email, sql.NullString{}, sql.NullString{}, sql.NullInt64{}, sql.NullString{}, sql.NullString{})
	hashPassword := HashPassword(password)
	user.Password = string(hashPassword)

	err := user.Save()
	return user.ID, err
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
		return nil, err
	}

	user := NewUser(id, email, provider, oauthToken, uuid, userName, avatarURL)
	bytePassword := []byte(userPassword)
	err = bcrypt.CompareHashAndPassword([]byte(password), bytePassword)
	if err != nil {
		return nil, err
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
		return nil, err
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
		panic(err)
	}
	for rows.Next() {
		err := rows.Scan(&id, &email, &provider, &oauthToken, &userName, &uuid, &avatarURL)
		if err != nil {
			panic(err)
		}
	}
	user := NewUser(id, email, provider, oauthToken, uuid, userName, avatarURL)

	// create
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

func (u *UserStruct) Projects() []*project.ProjectStruct {
	table := u.database.Init()
	defer table.Close()

	var slice []*project.ProjectStruct
	rows, err := table.Query("select id, user_id, repository_id, title, description, show_issues, show_pull_requests from projects where user_id = ?;", u.ID)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var id, userID int64
		var repositoryID sql.NullInt64
		var title string
		var description string
		var showIssues, showPullRequests bool
		err := rows.Scan(&id, &userID, &repositoryID, &title, &description, &showIssues, &showPullRequests)
		if err != nil {
			panic(err)
		}
		if id != 0 {
			p := project.NewProject(id, userID, title, description, repositoryID, showIssues, showPullRequests)
			slice = append(slice, p)
		}
	}
	return slice
}

// UniqueValidation is check record unique
func (u *UserStruct) UniqueValidation(tx *sql.Tx) bool {
	// unique
	var id int64
	var err error
	if u.ID == 0 {
		err = tx.QueryRow("select id from users where email = ?;", u.Email).Scan(&id)
	} else {
		err = tx.QueryRow("select id from users where email = ? and id != ?;", u.Email, u.ID).Scan(&id)
	}
	if err == nil || id != 0 {
		return false
	}
	return true
}

func (u *UserStruct) Save() error {
	table := u.database.Init()
	defer table.Close()
	tx, err := table.Begin()
	if err != nil {
		panic(err)
	}

	if !u.UniqueValidation(tx) {
		tx.Rollback()
		return errors.New("record is not unique")
	}
	result, err := tx.Exec("insert into users (email, password, provider, oauth_token, uuid, user_name, avatar_url, created_at) values (?, ?, ?, ?, ?, ?, ?, now());", u.Email, u.Password, u.Provider, u.OauthToken, u.Uuid, u.UserName, u.Avatar)
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	tx.Commit()
	u.ID, _ = result.LastInsertId()
	logging.SharedInstance().MethodInfo("user", "Save", false).Infof("user saved: %v", u.ID)
	return nil
}

func (u *UserStruct) Update() error {
	table := u.database.Init()
	defer table.Close()
	tx, err := table.Begin()
	if err != nil {
		panic(err)
	}

	if !u.UniqueValidation(tx) {
		tx.Rollback()
		return errors.New("record is not unique")
	}

	_, err = tx.Exec("update users set provider = ?, oauth_token = ?, uuid = ?, user_name = ?, avatar_url = ? where email = ?;", u.Provider, u.OauthToken, u.Uuid, u.UserName, u.Avatar, u.Email)
	if err != nil {
		panic(err)
	}
	tx.Commit()
	logging.SharedInstance().MethodInfo("user", "Update", false).Infof("user updated: %v", u.ID)
	return nil
}

func (u *UserStruct) CreateGithubUser(token string, githubUser *github.User, primaryEmail string) error {
	u.Email = primaryEmail
	bytePassword := HashPassword(randomString())

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
