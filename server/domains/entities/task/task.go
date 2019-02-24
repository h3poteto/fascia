package task

import (
	"database/sql"

	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/server/domains/entities/list_option"
	"github.com/h3poteto/fascia/server/domains/entities/repository"
	"github.com/h3poteto/fascia/server/infrastructures/task"
	"github.com/pkg/errors"
)

// Task has a task model object
type Task struct {
	ID             int64
	ListID         int64
	ProjectID      int64
	UserID         int64
	IssueNumber    sql.NullInt64
	Title          string
	Description    string
	PullRequest    bool
	HTMLURL        sql.NullString
	infrastructure *task.Task
}

// New returns a task entity
func New(id, listID, projectID, userID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString) *Task {
	infrastructure := task.New(id, listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL)
	t := &Task{
		infrastructure: infrastructure,
	}
	t.reload()
	return t
}

func (t *Task) reflect() {
	t.infrastructure.ID = t.ID
	t.infrastructure.ListID = t.ListID
	t.infrastructure.ProjectID = t.ProjectID
	t.infrastructure.UserID = t.UserID
	t.infrastructure.IssueNumber = t.IssueNumber
	t.infrastructure.Title = t.Title
	t.infrastructure.Description = t.Description
	t.infrastructure.PullRequest = t.PullRequest
	t.infrastructure.HTMLURL = t.HTMLURL
}

func (t *Task) reload() error {
	if t.ID != 0 {
		latestTask, err := task.Find(t.ID)
		if err != nil {
			return err
		}
		t.infrastructure = latestTask
	}
	t.ID = t.infrastructure.ID
	t.ListID = t.infrastructure.ListID
	t.ProjectID = t.infrastructure.ProjectID
	t.UserID = t.infrastructure.UserID
	t.IssueNumber = t.infrastructure.IssueNumber
	t.Title = t.infrastructure.Title
	t.Description = t.infrastructure.Description
	t.PullRequest = t.infrastructure.PullRequest
	t.HTMLURL = t.infrastructure.HTMLURL
	return nil
}

// Save call save in model
func (t *Task) Save() error {
	t.reflect()
	if err := t.infrastructure.Save(); err != nil {
		return err
	}
	return t.reload()
}

// Update call update in model
func (t *Task) Update(listID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString) error {
	err := t.infrastructure.Update(listID, issueNumber, title, description, pullRequest, htmlURL)
	if err != nil {
		return err
	}
	err = t.reload()
	if err != nil {
		return err
	}
	return nil
}

// ChangeList change list of a task, and reorder task
// returns isReorder, error.
func (t *Task) ChangeList(listID int64, prevToTaskID *int64) (bool, error) {
	var isReorder bool
	// リストを移動させるのか同リスト内の並び替えなのかどうかを見て，並び替えならgithub同期したくない
	if listID == t.infrastructure.ListID {
		isReorder = true
	} else {
		isReorder = false
	}

	err := t.infrastructure.ChangeList(listID, prevToTaskID)
	if err != nil {
		return isReorder, err
	}
	err = t.reload()
	if err != nil {
		return isReorder, err
	}
	return isReorder, nil
}

// Delete call delete in model
func (t *Task) Delete() error {
	err := t.infrastructure.Delete()
	if err != nil {
		return err
	}
	t.infrastructure = nil
	return nil
}

// SyncIssue apply task information to github issue, and take issue to task
func (t *Task) SyncIssue(repo *repository.Repository, token string) (*github.Issue, error) {
	listTitle, listColor, listOptionID, err := t.infrastructure.List()
	if err != nil {
		return nil, err
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
		issueAction = listOption.Action
	}

	var issue *github.Issue
	// issueを確認する
	if t.IssueNumber.Valid {
		issue, err = repo.GetGithubIssue(token, int(t.IssueNumber.Int64))
		if err != nil {
			return nil, err
		}
	}

	// issueがない場合には作成する
	if issue == nil {
		issue, err = repo.CreateGithubIssue(token, t.Title, t.Description, labelName)
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
	result, err := repo.EditGithubIssue(token, t.Title, t.Description, issueAction, *issue.Number, labelName)
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
