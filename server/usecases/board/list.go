package board

import (
	"database/sql"

	"github.com/h3poteto/fascia/lib/modules/database"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/domains/repo"
	"github.com/h3poteto/fascia/server/domains/task"
	repository "github.com/h3poteto/fascia/server/infrastructures/list"
)

// InjectDB set DB connection from connection pool.
func InjectDB() *sql.DB {
	return database.SharedInstance().Connection
}

// InjectListRepository returns list Repository.
func InjectListRepository() list.Repository {
	return repository.New(InjectDB())
}

// ListHasCloseAction returns either close action or not.
func ListHasCloseAction(list *list.List) (bool, error) {
	return list.HasCloseAction()
}

// FindList returns a list.
func FindList(projectID, listID int64) (*list.List, error) {
	repo := InjectListRepository()
	return repo.Find(projectID, listID)
}

// CreateList creates a list, and sync to github.
func CreateList(projectID, userID int64, title, color string, optionID sql.NullInt64, isHidden bool) (*list.List, error) {
	nullableTitle := sql.NullString{String: title, Valid: true}
	nullableColor := sql.NullString{String: color, Valid: true}
	repo := InjectListRepository()
	id, err := repo.Create(projectID, userID, nullableTitle, nullableColor, optionID, isHidden, nil)
	if err != nil {
		return nil, err
	}
	l, err := repo.Find(projectID, id)
	if err != nil {
		return nil, err
	}

	go func(l *list.List) {
		projectID := l.ProjectID
		p, err := FindProject(projectID)
		// TODO: log
		if err != nil {
			return
		}
		token, err := p.OauthToken()
		if err != nil {
			return
		}
		repo, err := p.Repository(InjectRepoRepository())
		if err != nil {
			return
		}
		err = fetchCreatedList(l, token, repo)
		if err != nil {
			return
		}
	}(l)
	return l, nil
}

func fetchCreatedList(l *list.List, oauthToken string, repo *repo.Repo) error {
	if repo != nil {
		label, err := repo.CheckLabelPresent(oauthToken, l.Title.String)
		if err != nil {
			return err
		} else if label == nil {
			label, err = repo.CreateGithubLabel(oauthToken, l.Title.String, l.Color.String)
			if err != nil {
				return err
			}
		} else {
			// 色だけはこちら指定のものに変更したい
			_, err := repo.UpdateGithubLabel(oauthToken, l.Title.String, l.Title.String, l.Color.String)
			if err != nil {
				return err
			}
			logging.SharedInstance().MethodInfo("list", "Save").Info("github label already exist")
		}
	}
	return nil
}

// UpdateList updates a list, and sync to github.
func UpdateList(l *list.List, title, color string, optionID int64) (*list.List, error) {
	nullableTitle := sql.NullString{String: title, Valid: true}
	nullableColor := sql.NullString{String: color, Valid: true}
	repo := InjectListRepository()
	// Allow null option, so ignore error.
	option, _ := repo.FindOptionByID(optionID)
	err := l.Update(nullableTitle, nullableColor, option)
	if err != nil {
		return nil, err
	}
	if err = repo.Update(l); err != nil {
		return nil, err
	}

	go func(l *list.List, title, color string) {
		projectID := l.ProjectID
		p, err := FindProject(projectID)
		if err != nil {
			return
		}
		token, err := p.OauthToken()
		if err != nil {
			return
		}
		repo, err := p.Repository(InjectRepoRepository())
		if err != nil {
			return
		}
		err = fetchUpdatedList(l, token, repo, title, color)
		if err != nil {
			return
		}
	}(l, title, color)
	return l, nil
}

func fetchUpdatedList(l *list.List, oauthToken string, repo *repo.Repo, newTitle, newColor string) error {
	if repo != nil {
		// 編集前のラベルがそもそも存在しているかどうかを確認する
		existLabel, err := repo.CheckLabelPresent(oauthToken, l.Title.String)
		if err != nil {
			return err
		} else if existLabel == nil {
			// editの場合ここに入る可能性は，createのgithub同期がうまく動いていない場合のみ
			// 編集前のラベルが存在しなければ新しく作る
			_, err := repo.CreateGithubLabel(oauthToken, newTitle, newColor)
			if err != nil {
				return err
			}
		} else {
			_, err := repo.UpdateGithubLabel(oauthToken, l.Title.String, newTitle, newColor)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func listsWithCloseAction(lists []list.List) []list.List {
	var closeLists []list.List
	for _, list := range lists {
		result, err := list.HasCloseAction()
		if err != nil {
			logging.SharedInstance().MethodInfo("Project", "listsWithCloseAction").Info(err)
		} else if result {
			closeLists = append(closeLists, list)
		}
	}
	return closeLists
}

// ListTasks returns tasks related the list.
func ListTasks(l *list.List) ([]*task.Task, error) {
	repo := InjectTaskRepository()
	return task.Tasks(l.ID, repo)
}

// ListOptionAll returns all list options
func ListOptionAll() ([]*list.Option, error) {
	repo := InjectListRepository()
	return repo.AllOption()
}

// FindListOptionByID returns a list option service
func FindListOptionByID(id int64) (*list.Option, error) {
	repo := InjectListRepository()
	return repo.FindOptionByID(id)
}

// FindListOptionByAction returns a list option service
func FindListOptionByAction(action string) (*list.Option, error) {
	repo := InjectListRepository()
	return repo.FindOptionByAction(action)
}

// HideList hides a list.
func HideList(l *list.List) error {
	l.Hide()
	repo := InjectListRepository()
	return repo.Update(l)
}

// DisplayList display a list.
func DisplayList(l *list.List) error {
	l.Display()
	repo := InjectListRepository()
	return repo.Update(l)
}
