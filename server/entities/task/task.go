package task

import (
	"database/sql"

	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/database"
	"github.com/h3poteto/fascia/server/entities/list_option"
	"github.com/h3poteto/fascia/server/entities/repository"
	"github.com/h3poteto/fascia/server/models/task"
	"github.com/pkg/errors"
)

// Task has a task model object
type Task struct {
	TaskModel *task.Task
	db        *sql.DB
}

// New returns a task entity
func New(id, listID, projectID, userID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString) *Task {
	return &Task{
		TaskModel: task.New(id, listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL),
		db:        database.SharedInstance().Connection,
	}
}

// Find returns a task entity
func Find(listID, taskID int64) (*Task, error) {
	t, err := task.Find(listID, taskID)
	if err != nil {
		return nil, err
	}
	return &Task{
		TaskModel: t,
		db:        database.SharedInstance().Connection,
	}, nil
}

// FindByIssueNumber returns a task entity
func FindByIssueNumber(projectID int64, issueNumber int) (*Task, error) {
	t, err := task.FindByIssueNumber(projectID, issueNumber)
	if err != nil {
		return nil, err
	}
	return &Task{
		TaskModel: t,
		db:        database.SharedInstance().Connection,
	}, nil
}

// Save call save in model
func (t *Task) Save() error {
	return t.TaskModel.Save()
}

// Update call update in model
func (t *Task) Update(listID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString) error {
	return t.TaskModel.Update(listID, issueNumber, title, description, pullRequest, htmlURL)
}

// ChangeList change list of a task, and reorder task
// returns isReorder, error.
func (t *Task) ChangeList(listID int64, prevToTaskID *int64) (bool, error) {
	var isReorder bool
	// リストを移動させるのか同リスト内の並び替えなのかどうかを見て，並び替えならgithub同期したくない
	if listID == t.TaskModel.ListID {
		isReorder = true
	} else {
		isReorder = false
	}

	return isReorder, t.TaskModel.ChangeList(listID, prevToTaskID)
}

// Delete call delete in model
func (t *Task) Delete() error {
	return t.TaskModel.Delete()
}

// SyncIssue apply task information to github issue, and take issue to task
func (t *Task) SyncIssue(repo *repository.Repository, token string) (*github.Issue, error) {
	var listTitle, listColor sql.NullString
	var listOptionID sql.NullInt64
	err := t.db.QueryRow("select title, color, list_option_id from lists where id = ?;", t.TaskModel.ListID).Scan(&listTitle, &listColor, &listOptionID)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}

	labelName, err := t.syncLabel(listTitle.String, listColor.String, token, repo)
	if err != nil {
		return nil, err
	}

	// list_optionが定義しているactionを先に出しておく
	var issueAction string
	var listOption *list_option.ListOption
	if listOptionID.Valid {
		listOption, err := list_option.FindByID(listOptionID.Int64)
		if err != nil {
			return nil, err
		}
		issueAction = listOption.ListOptionModel.Action
	}

	var issue *github.Issue
	// issueを確認する
	if t.TaskModel.IssueNumber.Valid {
		issue, err = repo.GetGithubIssue(token, int(t.TaskModel.IssueNumber.Int64))
		if err != nil {
			return nil, err
		}
	}

	// issueがない場合には作成する
	if issue == nil {
		issue, err = repo.CreateGithubIssue(token, t.TaskModel.Title, t.TaskModel.Description, labelName)
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
	result, err := repo.EditGithubIssue(token, t.TaskModel.Title, t.TaskModel.Description, issueAction, *issue.Number, labelName)
	if err != nil {
		return nil, err
	}
	if !result {
		return nil, errors.New("unexpected error")
	}
	return issue, nil
}

func (t *Task) syncLabel(listTitle string, listColor string, token string, repo *repository.Repository) ([]string, error) {
	if listTitle == config.Element("init_list").(map[interface{}]interface{})["none"].(string) {
		return []string{}, nil
	}
	label, err := repo.CheckLabelPresent(token, listTitle)
	if err != nil {
		return nil, err
	} else if label == nil {
		// 対象のラベルがない場合には新規作成する
		label, err = repo.CreateGithubLabel(token, listTitle, listColor)
		if label == nil {
			return nil, err
		}
	}
	return []string{*label.Name}, nil
}
