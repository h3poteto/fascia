package services

import (
	"database/sql"

	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/domains/repo"
	"github.com/h3poteto/fascia/server/domains/task"
	"github.com/pkg/errors"
)

// AfterCreateTask fetch the created task.
func AfterCreateTask(t *task.Task, projectInfra project.Repository, listInfra list.Repository, taskInfra task.Repository, repoInfra repo.Repository) {
	projectID := t.ProjectID
	p, err := projectInfra.Find(projectID)
	// TODO: log
	if err != nil {
		return
	}
	token, err := projectInfra.OauthToken(p.ID)
	if err != nil {
		return
	}
	repo, err := repoInfra.FindByProjectID(p.ID)
	if err != nil {
		return
	}
	err = fetchCreatedTask(t, token, repo, listInfra, taskInfra)
	if err != nil {
		return
	}
}

func fetchCreatedTask(t *task.Task, oauthToken string, repo *repo.Repo, listInfra list.Repository, taskInfra task.Repository) error {
	if repo != nil {
		issue, err := syncTaskToIssue(t, repo, oauthToken, listInfra)
		if err != nil {
			return errors.Wrap(err, "sync github error")
		}
		issueNumber := sql.NullInt64{Int64: int64(*issue.Number), Valid: true}
		HTMLURL := sql.NullString{String: *issue.HTMLURL, Valid: true}

		t.Update(
			t.ListID,
			issueNumber,
			t.Title,
			t.Description,
			t.PullRequest,
			HTMLURL,
		)
		err = taskInfra.Update(
			t.ID,
			t.ListID,
			t.ProjectID,
			t.UserID,
			t.IssueNumber,
			t.Title,
			t.Description,
			t.PullRequest,
			t.HTMLURL,
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

// syncTaskToIssue apply task information to github issue, and take issue to task
func syncTaskToIssue(t *task.Task, r *repo.Repo, token string, listInfra list.Repository) (*github.Issue, error) {
	l, err := listInfra.Find(t.ProjectID, t.ListID)
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
	if l.Option != nil {
		listOption, err := listInfra.FindOptionByID(l.Option.ID)
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

func syncLabel(t *task.Task, listTitle string, listColor string, token string, r *repo.Repo) ([]string, error) {
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
