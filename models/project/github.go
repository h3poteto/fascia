package project

import (
	"../../config"
	"../../modules/hub"
	"../../modules/logging"
	"../db"
	"../list"
	"../task"
	"database/sql"
	"errors"

	"github.com/google/go-github/github"
)

// LoadFromGithub load tasks from github issues
func (u *ProjectStruct) LoadFromGithub(issues []github.Issue) error {
	for _, issue := range issues {
		targetTask, _ := task.FindByIssueNumber(u.ID, *issue.Number)

		err := u.taskApplyLabel(targetTask, &issue)
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
		return err
	}
	parentProject, err := FindProject(projectID)
	if err != nil {
		return err
	}
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
		return err
	}
	parentProject, err := FindProject(projectID)
	if err != nil {
		return err
	}
	targetTask, _ := task.FindByIssueNumber(projectID, *body.Number)

	// TODO: もしgithubへのアクセスが増大するようであれば，PullRequestオブジェクトからラベルの付替えを行うように改修する

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

// GithubLabels get issues which match project's lists from github
func GithubLabels(issue *github.Issue, projectLists []*list.ListStruct) []list.ListStruct {
	var githubLabels []list.ListStruct
	for _, label := range issue.Labels {
		for _, list := range projectLists {
			if list.Title.Valid && list.Title.String == *label.Name {
				githubLabels = append(githubLabels, *list)
			}
		}
	}
	return githubLabels
}

func (u *ProjectStruct) taskApplyLabel(targetTask *task.TaskStruct, issue *github.Issue) error {
	if targetTask == nil {
		err := u.createNewTask(issue)
		if err != nil {
			logging.SharedInstance().MethodInfo("Project", "taskApplyLabel", true).Errorf("create new task failed: %v", err)
			return err
		}
		return nil
	}
	issueTask, err := u.applyListToTask(targetTask, issue)
	if err != nil {
		logging.SharedInstance().MethodInfo("Project", "taskApplyLabel", true).Errorf("apply list to task failed: %v", err)
		return err
	}
	if err := issueTask.Update(nil, nil); err != nil {
		logging.SharedInstance().MethodInfo("Project", "taskApplyLabel", true).Error("task update failed")
		return err
	}
	return nil
}

func (u *ProjectStruct) reopenTask(targetTask *task.TaskStruct, issue *github.Issue) error {
	issueTask, err := u.applyListToTask(targetTask, issue)
	if err != nil {
		logging.SharedInstance().MethodInfo("Project", "reopenTask", true).Errorf("apply list to task failed: %v", err)
		return err
	}
	if err := issueTask.Update(nil, nil); err != nil {
		logging.SharedInstance().MethodInfo("Project", "reopenTask", true).Error("task update failed")
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
		logging.SharedInstance().MethodInfo("Project", "createNewTask", true).Errorf("apply list to task failed: %v", err)
		return err
	}
	if err := issueTask.Save(nil, nil); err != nil {
		logging.SharedInstance().MethodInfo("Project", "createNewTask", true).Error("task save failed")
		return err
	}
	return nil

}

func (u *ProjectStruct) applyListToTask(issueTask *task.TaskStruct, issue *github.Issue) (*task.TaskStruct, error) {
	// close noneの用意
	var closedList, noneList *list.ListStruct
	for _, list := range u.Lists() {
		if list.Title.Valid && list.Title.String == config.Element("init_list").(map[interface{}]interface{})["done"].(string) {
			closedList = list
		}
	}
	if closedList == nil {
		logging.SharedInstance().MethodInfo("Project", "applyListToTask", true).Panic("cannot find close list")
		return nil, errors.New("cannot find close list")
	}

	noneList, err := u.NoneList()
	if err != nil {
		logging.SharedInstance().MethodInfo("Project", "applyListToTask", true).Panic("cannot find none list")
		return nil, err
	}

	githubLabels := GithubLabels(issue, u.Lists())

	// label所属よりcloseかどうかを優先して判定したい
	// closeのものはどんなlabelがついていようと，doneに放り込む
	if *issue.State == "open" {
		if len(githubLabels) == 1 {
			// 一つのlistだけが該当するとき
			issueTask.ListID = githubLabels[0].ID
		} else if len(githubLabels) > 1 {
			// 複数のlistが該当するとき
			issueTask.ListID = githubLabels[0].ID
		} else {
			// listに該当しないlabelしか持っていない
			// or そもそもlabelがひとつもついていない
			issueTask.ListID = noneList.ID
		}
	} else {
		issueTask.ListID = closedList.ID
	}
	return issueTask, nil
}
