package board

import (
	"database/sql"

	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/domains/repo"
	domain "github.com/h3poteto/fascia/server/domains/task"
	repository "github.com/h3poteto/fascia/server/infrastructures/task"
	"github.com/pkg/errors"
)

// InjectTaskRepository returns a task Repository.
func InjectTaskRepository() domain.Repository {
	return repository.New(InjectDB())
}

// FindTask finds a task.
func FindTask(id int64) (*domain.Task, error) {
	return domain.Find(id, InjectTaskRepository())
}

// CreateTask creates a task, and sync to github.
func CreateTask(listID, projectID, userID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString) (*domain.Task, error) {
	task := domain.New(0, listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL, InjectTaskRepository())
	err := task.Create()
	if err != nil {
		return nil, err
	}

	go func(task *domain.Task) {
		projectID := task.ProjectID
		p, err := FindProject(projectID)
		// TODO: log
		if err != nil {
			return
		}
		infra := InjectProjectRepository()
		token, err := infra.OauthToken(p.ID)
		if err != nil {
			return
		}
		repo, err := ProjectRepository(p)
		if err != nil {
			return
		}
		err = fetchCreatedTask(task, token, repo)
		if err != nil {
			return
		}
	}(task)

	return task, nil
}

// UpdateTask updates a task, and sync to github.
func UpdateTask(task *domain.Task, listID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString) error {
	err := task.Update(listID, issueNumber, title, description, pullRequest, htmlURL)
	if err != nil {
		return err
	}

	go func(task *domain.Task) {
		projectID := task.ProjectID
		p, err := FindProject(projectID)
		// TODO: log
		if err != nil {
			return
		}
		infra := InjectProjectRepository()
		token, err := infra.OauthToken(p.ID)
		if err != nil {
			return
		}
		repo, err := ProjectRepository(p)
		if err != nil {
			return
		}
		err = fetchUpdatedTask(task, token, repo)
		if err != nil {
			return
		}
	}(task)

	return err
}

// TaskChangeList changes the task and sync it to github.
func TaskChangeList(task *domain.Task, listID int64, prevToTaskID *int64) error {
	isReorder, err := task.ChangeList(listID, prevToTaskID)
	if err != nil {
		return err
	}

	go func(task *domain.Task, isReorder bool) {
		projectID := task.ProjectID
		p, err := FindProject(projectID)
		// TODO: log
		if err != nil {
			return
		}
		infra := InjectProjectRepository()
		token, err := infra.OauthToken(p.ID)
		if err != nil {
			return
		}
		repo, err := ProjectRepository(p)
		if err != nil {
			return
		}
		err = fetchChangedList(task, token, repo, isReorder)
		if err != nil {
			return
		}
	}(task, isReorder)
	return nil
}

func fetchCreatedTask(t *domain.Task, oauthToken string, repo *repo.Repo) error {
	if repo != nil {
		issue, err := syncTaskToIssue(t, repo, oauthToken)
		if err != nil {
			return errors.Wrap(err, "sync github error")
		}
		issueNumber := sql.NullInt64{Int64: int64(*issue.Number), Valid: true}
		HTMLURL := sql.NullString{String: *issue.HTMLURL, Valid: true}

		err = t.Update(
			t.ListID,
			issueNumber,
			t.Title,
			t.Description,
			t.PullRequest,
			HTMLURL,
		)
		if err != nil {
			// note: この時にはすでにissueが作られてしまっているが，DBへの保存には失敗したということ
			// 本来であれば，issueを削除しなければならない
			// しかし，githubにはissue削除APIがない
			// 運がいいことに，Webhookが正常に動作していれば，作られたissueに応じてDBにタスクを作ってくれる
			// そちらに任せることにしよう
			return errors.Wrap(err, "sql execute error")
		}
		logging.SharedInstance().MethodInfo("task", "Save").Info("issue number is updated")
	}
	return nil
}

func fetchUpdatedTask(t *domain.Task, oauthToken string, repo *repo.Repo) error {
	// github側へ同期
	if repo != nil {
		_, err := syncTaskToIssue(t, repo, oauthToken)
		if err != nil {
			logging.SharedInstance().MethodInfo("task", "Update").Error(err)
			return err
		}
		logging.SharedInstance().MethodInfo("task", "Update").Debugf("task synced to github: %+v", t)
	}
	return nil
}

func fetchChangedList(t *domain.Task, oauthToken string, repo *repo.Repo, isReorder bool) error {
	if !isReorder && repo != nil {
		_, err := syncTaskToIssue(t, repo, oauthToken)
		if err != nil {
			logging.SharedInstance().MethodInfo("task", "ChangeList").Error(err)
			return err
		}
		logging.SharedInstance().MethodInfo("Task", "Update").Debugf("task synced to github: %+v", t)
	}
	return nil
}

// syncTaskToIssue apply task information to github issue, and take issue to task
func syncTaskToIssue(t *domain.Task, r *repo.Repo, token string) (*github.Issue, error) {
	l, err := FindList(t.ProjectID, t.ListID)
	if err != nil {
		return nil, err
	}

	labelName, err := syncLabel(t, l.Title.String, l.Color.String, token, r)
	if err != nil {
		return nil, err
	}

	// Get all actions which are defined at list option.
	var issueAction string
	var listOption *list.Option
	listRepo := InjectListRepository()
	if l.Option != nil {
		listOption, err := listRepo.FindOptionByID(l.Option.ID)
		if err != nil {
			return nil, err
		}
		issueAction = listOption.Action
	}

	var issue *github.Issue
	// Confirm the issue.
	if t.IssueNumber.Valid {
		issue, err = r.GetGithubIssue(token, int(t.IssueNumber.Int64))
		if err != nil {
			return nil, err
		}
	}

	// Create an issue when the issue does not exist.
	if issue == nil {
		issue, err = r.CreateGithubIssue(token, t.Title, t.Description, labelName)
		if err != nil {
			return nil, err
		}
		// もしlist_optionがopenだった場合には，これ以上更新する必要がない
		// むしろこれ以上処理を続けさせると，Createした直後にUpdateがかかってしまい，github側に編集履歴が残ってしまう
		// そのため，ここで終わりにする
		// listOptionは特に指定しない場合にはnilになっているので，nilのパターンも除外しなければならない
		if listOption == nil || !listOption.IsCloseAction() {
			return issue, nil
		}
	}

	// issueがある場合には更新する
	// このときにlist_optionが定義するaction通りにissueを更新する
	result, err := r.EditGithubIssue(token, t.Title, t.Description, issueAction, *issue.Number, labelName)
	if err != nil {
		return nil, err
	}
	if !result {
		return nil, errors.New("unexpected error")
	}
	return issue, nil
}

func syncLabel(t *domain.Task, listTitle string, listColor string, token string, r *repo.Repo) ([]string, error) {
	if listTitle == config.Element("init_list").(map[interface{}]interface{})["none"].(string) {
		return []string{}, nil
	}
	label, err := r.CheckLabelPresent(token, listTitle)
	if err != nil {
		return nil, err
	} else if label == nil {
		// 対象のラベルがない場合には新規作成する
		label, err = r.CreateGithubLabel(token, listTitle, listColor)
		if label == nil {
			return nil, err
		}
	}
	return []string{*label.Name}, nil
}
