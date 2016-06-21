package project

import (
	"../../config"
	"../../modules/hub"
	"../../modules/logging"
	"../db"
	"../list"
	"../task"

	"database/sql"
	"strings"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

// ListLoadFromGithub load lists from github labels
func (u *ProjectStruct) ListLoadFromGithub(labels []*github.Label) error {
	lists, err := u.Lists()
	if err != nil {
		return err
	}
	for _, l := range lists {
		if err := u.labelUpdate(l, labels); err != nil {
			return err
		}
	}
	return nil
}

func (u *ProjectStruct) labelUpdate(l *list.ListStruct, labels []*github.Label) error {
	for _, label := range labels {
		if strings.ToLower(*label.Name) == strings.ToLower(l.Title.String) {
			l.Color.String = *label.Color
			if err := l.UpdateColor(); err != nil {
				return err
			}
		}
	}
	return nil
}

// TaskLoadFromGithub load tasks from github issues
func (u *ProjectStruct) TaskLoadFromGithub(issues []*github.Issue) error {
	for _, issue := range issues {
		targetTask, _ := task.FindByIssueNumber(u.ID, *issue.Number)

		err := u.taskApplyLabel(targetTask, issue)
		if err != nil {
			return err
		}
	}
	return nil
}

// IssuesEvent apply issue changes to task
func IssuesEvent(repositoryID int64, body github.IssuesEvent) error {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var projectID int64
	err := table.QueryRow("select id from projects where repository_id = ?;", repositoryID).Scan(&projectID)
	if err != nil {
		return errors.Wrap(err, "sql select error")
	}
	parentProject, err := FindProject(projectID)
	if err != nil {
		return err
	}
	// taskが見つからない場合は新規作成するのでエラーハンドリング不要
	targetTask, _ := task.FindByIssueNumber(projectID, *body.Issue.Number)

	switch *body.Action {
	case "opened", "reopened":
		if targetTask == nil {
			err = parentProject.createNewTask(body.Issue)
		} else {
			err = parentProject.reopenTask(targetTask, body.Issue)
		}
	case "closed", "labeled", "unlabeled":
		err = parentProject.taskApplyLabel(targetTask, body.Issue)
	}
	return err
}

// PullRequestEvent apply issue changes to task
func PullRequestEvent(repositoryID int64, body github.PullRequestEvent) error {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var projectID int64
	err := table.QueryRow("select id from projects where repository_id = ?;", repositoryID).Scan(&projectID)
	if err != nil {
		return errors.Wrap(err, "sql select error")
	}
	parentProject, err := FindProject(projectID)
	if err != nil {
		return err
	}
	// taskが見つからない場合は新規作成するのでエラーハンドリング不要
	targetTask, _ := task.FindByIssueNumber(projectID, *body.Number)

	// note: もしgithubへのアクセスが増大するようであれば，PullRequestオブジェクトからラベルの付替えを行うように改修する

	oauthToken, err := parentProject.OauthToken()
	if err != nil {
		return err
	}

	repo, err := parentProject.Repository()
	if err != nil {
		return err
	}
	issue, err := hub.GetGithubIssue(oauthToken, repo, *body.Number)

	switch *body.Action {
	case "opened", "reopened":
		if targetTask == nil {
			err = parentProject.createNewTask(issue)
		} else {
			err = parentProject.reopenTask(targetTask, issue)
		}
	case "closed", "labeled", "unlabeled":
		err = parentProject.taskApplyLabel(targetTask, issue)
	}
	return err
}

// githubLabels get issues which match project's lists from github
func githubLabelLists(issue *github.Issue, projectLists []*list.ListStruct) []list.ListStruct {
	var githubLabels []list.ListStruct
	for _, label := range issue.Labels {
		for _, list := range projectLists {
			if list.Title.Valid && strings.ToLower(list.Title.String) == strings.ToLower(*label.Name) {
				githubLabels = append(githubLabels, *list)
			}
		}
	}
	return githubLabels
}

func listsWithCloseAction(lists []list.ListStruct) []list.ListStruct {
	var closeLists []list.ListStruct
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

func (u *ProjectStruct) taskApplyLabel(targetTask *task.TaskStruct, issue *github.Issue) error {
	if targetTask == nil {
		err := u.createNewTask(issue)
		if err != nil {
			return err
		}
		return nil
	}
	issueTask, err := u.applyListToTask(targetTask, issue)
	if err != nil {
		return err
	}
	issueTask, err = u.applyIssueInfoToTask(issueTask, issue)
	if err != nil {
		return err
	}
	if err := issueTask.Update(nil, nil); err != nil {
		return err
	}
	return nil
}

func (u *ProjectStruct) reopenTask(targetTask *task.TaskStruct, issue *github.Issue) error {
	issueTask, err := u.applyListToTask(targetTask, issue)
	if err != nil {
		return err
	}
	if err := issueTask.Update(nil, nil); err != nil {
		return err
	}
	return nil
}

func (u *ProjectStruct) createNewTask(issue *github.Issue) error {

	issueTask := task.NewTask(
		0,
		0,
		u.ID,
		u.UserID,
		sql.NullInt64{Int64: int64(*issue.Number), Valid: true},
		*issue.Title,
		*issue.Body,
		hub.IsPullRequest(issue),
		sql.NullString{String: *issue.HTMLURL, Valid: true},
	)

	issueTask, err := u.applyListToTask(issueTask, issue)
	if err != nil {
		return err
	}
	if err := issueTask.Save(nil, nil); err != nil {
		return err
	}
	return nil

}

func (u *ProjectStruct) applyListToTask(issueTask *task.TaskStruct, issue *github.Issue) (*task.TaskStruct, error) {
	// close noneの用意
	var closedList, noneList *list.ListStruct
	lists, err := u.Lists()
	if err != nil {
		return nil, err
	}
	for _, list := range lists {
		if list.Title.Valid && list.Title.String == config.Element("init_list").(map[interface{}]interface{})["done"].(string) {
			closedList = list
		}
	}
	if closedList == nil {
		return nil, errors.New("cannot find close list")
	}

	noneList, err = u.NoneList()
	if err != nil {
		return nil, err
	}

	labelLists := githubLabelLists(issue, lists)
	logging.SharedInstance().MethodInfo("Project", "applyListToTask").Debugf("github label: %+v", labelLists)

	// label所属よりcloseかどうかを優先して判定したい
	if *issue.State == "open" {
		if len(labelLists) >= 1 {
			// 一以上listだけが該当するとき
			issueTask.ListID = labelLists[0].ID
		} else {
			// listに該当しないlabelしか持っていない
			// or そもそもlabelがひとつもついていない
			issueTask.ListID = noneList.ID
		}
	} else {
		// closeのものは，該当するlistがあったとしてもそのまま放り込めない
		// Because: ToDoがついたままcloseされることはよくある
		// 該当するlistのlist_optionがcloseのときのみ，そのlistに放り込める
		listsWithClose := listsWithCloseAction(labelLists)
		logging.SharedInstance().MethodInfo("Project", "applyListToTask").Debugf("lists with close action: %+v", listsWithClose)
		if len(listsWithClose) >= 1 {
			issueTask.ListID = listsWithClose[0].ID
		} else {
			issueTask.ListID = closedList.ID
		}
	}
	return issueTask, nil
}

func (u *ProjectStruct) applyIssueInfoToTask(targetTask *task.TaskStruct, issue *github.Issue) (*task.TaskStruct, error) {
	if targetTask == nil {
		return nil, errors.New("target task is required")
	}
	targetTask.Title = *issue.Title
	targetTask.Description = *issue.Body
	targetTask.HTMLURL = sql.NullString{String: *issue.HTMLURL, Valid: true}
	return targetTask, nil
}
