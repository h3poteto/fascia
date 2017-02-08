package services

import (
	"database/sql"

	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/entities/repository"
	"github.com/h3poteto/fascia/server/entities/task"
	"github.com/pkg/errors"
)

// Task has a task entity
type Task struct {
	TaskEntity *task.Task
}

// NewTask returns a task service
func NewTask(id, listID, projectID, userID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString) *Task {
	return &Task{
		TaskEntity: task.New(id, listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL),
	}
}

// FindTask search task and returns a task service
func FindTask(listID, taskID int64) (*Task, error) {
	t, err := task.Find(listID, taskID)
	if err != nil {
		return nil, err
	}
	return &Task{
		TaskEntity: t,
	}, nil
}

// Save save a task, and fetch task to github
func (t *Task) Save() error {
	err := t.TaskEntity.Save()
	if err != nil {
		return err
	}

	go func(task *Task) {
		projectID := task.TaskEntity.TaskModel.ProjectID
		p, err := FindProject(projectID)
		// TODO: log
		if err != nil {
			return
		}
		token, err := p.ProjectEntity.OauthToken()
		if err != nil {
			return
		}
		repo, find, err := p.ProjectEntity.Repository()
		if err != nil {
			return
		}
		if !find {
			return
		}
		err = task.fetchCreated(token, repo)
		if err != nil {
			return
		}
	}(t)

	return nil
}

// Update update a task, and fetch task to github
func (t *Task) Update(listID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString) error {
	err := t.TaskEntity.Update(listID, issueNumber, title, description, pullRequest, htmlURL)
	if err != nil {
		return err
	}

	go func(task *Task) {
		projectID := task.TaskEntity.TaskModel.ProjectID
		p, err := FindProject(projectID)
		// TODO: log
		if err != nil {
			return
		}
		token, err := p.ProjectEntity.OauthToken()
		if err != nil {
			return
		}
		repo, find, err := p.ProjectEntity.Repository()
		if err != nil {
			return
		}
		if !find {
			return
		}
		err = task.fetchUpdated(token, repo)
		if err != nil {
			return
		}
	}(t)
	return nil
}

// ChangeList change list which task belongs, and fetch github
func (t *Task) ChangeList(listID int64, prevToTaskID *int64) error {
	isReorder, err := t.TaskEntity.ChangeList(listID, prevToTaskID)
	if err != nil {
		return err
	}

	go func(task *Task, isReorder bool) {
		projectID := task.TaskEntity.TaskModel.ProjectID
		p, err := FindProject(projectID)
		// TODO: log
		if err != nil {
			return
		}
		token, err := p.ProjectEntity.OauthToken()
		if err != nil {
			return
		}
		repo, find, err := p.ProjectEntity.Repository()
		if err != nil {
			return
		}
		if !find {
			return
		}
		err = task.fetchChangedList(token, repo, isReorder)
		if err != nil {
			return
		}
	}(t, isReorder)
	return nil
}

// Delete delete a task
func (t *Task) Delete() error {
	return t.TaskEntity.Delete()
}

func (t *Task) fetchCreated(oauthToken string, repo *repository.Repository) error {
	if repo != nil {
		issue, err := t.TaskEntity.SyncIssue(repo, oauthToken)
		if err != nil {
			return errors.Wrap(err, "sync github error")
		}
		issueNumber := sql.NullInt64{Int64: int64(*issue.Number), Valid: true}
		HTMLURL := sql.NullString{String: *issue.HTMLURL, Valid: true}

		err = t.TaskEntity.Update(
			t.TaskEntity.TaskModel.ListID,
			issueNumber,
			t.TaskEntity.TaskModel.Title,
			t.TaskEntity.TaskModel.Description,
			t.TaskEntity.TaskModel.PullRequest,
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

func (t *Task) fetchUpdated(oauthToken string, repo *repository.Repository) error {
	// github側へ同期
	if repo != nil {
		_, err := t.TaskEntity.SyncIssue(repo, oauthToken)
		if err != nil {
			logging.SharedInstance().MethodInfo("task", "Update").Error(err)
			return err
		}
		logging.SharedInstance().MethodInfo("task", "Update").Debugf("task synced to github: %+v", t)
	}
	return nil
}

func (t *Task) fetchChangedList(oauthToken string, repo *repository.Repository, isReorder bool) error {
	if !isReorder && repo != nil {
		_, err := t.TaskEntity.SyncIssue(repo, oauthToken)
		if err != nil {
			logging.SharedInstance().MethodInfo("task", "ChangeList").Error(err)
			return err
		}
		logging.SharedInstance().MethodInfo("Task", "Update").Debugf("task synced to github: %+v", t)
	}
	return nil
}
