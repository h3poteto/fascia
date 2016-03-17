package hook

import (
	"../../config"
	"../../modules/hub"
	"../../modules/logging"
	"../db"
	"../list"
	"../project"
	"../task"
	"database/sql"
	"errors"

	"github.com/google/go-github/github"
)

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
	parentProject := project.FindProject(projectID)
	targetTask, _ := task.FindByIssueNumber(projectID, *body.Issue.Number)

	switch *body.Action {
	case "opened", "reopened":
		if targetTask == nil {
			err = createNewTask(parentProject, body.Issue)
		} else {
			err = reopenTask(parentProject, targetTask, body.Issue)
		}
	case "closed", "labeled", "unlabeled":
		err = taskApplyLabel(parentProject, targetTask, body.Issue)
	}
	return err
}

func taskApplyLabel(parentProject *project.ProjectStruct, targetTask *task.TaskStruct, issue *github.Issue) error {
	if targetTask == nil {
		err := createNewTask(parentProject, issue)
		if err != nil {
			logging.SharedInstance().MethodInfo("Hook", "taskApplyLabel", true).Errorf("create new task failed: %v", err)
			return err
		}
		return nil
	}
	issueTask, err := applyListToTask(parentProject, targetTask, issue)
	if err != nil {
		logging.SharedInstance().MethodInfo("Hook", "taskApplyLabel", true).Errorf("apply list to task failed: %v", err)
		return err
	}
	if !issueTask.Update(nil, nil) {
		logging.SharedInstance().MethodInfo("Hook", "taskApplyLabel", true).Error("task update failed")
		return errors.New("update failed")
	}
	return nil
}

func reopenTask(parentProject *project.ProjectStruct, targetTask *task.TaskStruct, issue *github.Issue) error {
	issueTask, err := applyListToTask(parentProject, targetTask, issue)
	if err != nil {
		logging.SharedInstance().MethodInfo("Hook", "reopenTask", true).Errorf("apply list to task failed: %v", err)
		return err
	}
	if !issueTask.Update(nil, nil) {
		logging.SharedInstance().MethodInfo("Hook", "reopenTask", true).Error("task update failed")
		return errors.New("update failed")
	}
	return nil
}

func createNewTask(parentProject *project.ProjectStruct, issue *github.Issue) error {

	issueTask := task.NewTask(
		0,
		0,
		parentProject.ID,
		parentProject.UserID,
		sql.NullInt64{Int64: int64(*issue.Number), Valid: true},
		*issue.Title,
		*issue.Body,
		hub.IsPullRequest(issue),
		sql.NullString{String: *issue.HTMLURL, Valid: true},
	)

	issueTask, err := applyListToTask(parentProject, issueTask, issue)
	if err != nil {
		logging.SharedInstance().MethodInfo("Hook", "createNewTask", true).Errorf("apply list to task failed: %v", err)
		return err
	}
	if !issueTask.Save(nil, nil) {
		logging.SharedInstance().MethodInfo("Hook", "createNewTask", true).Error("task save failed")
	}
	return nil

}

func applyListToTask(parentProject *project.ProjectStruct, issueTask *task.TaskStruct, issue *github.Issue) (*task.TaskStruct, error) {
	// close noneの用意
	var closedList, noneList *list.ListStruct
	for _, list := range parentProject.Lists() {
		if list.Title.Valid && list.Title.String == config.Element("init_list").(map[interface{}]interface{})["done"].(string) {
			closedList = list
		}
	}
	if closedList == nil {
		logging.SharedInstance().MethodInfo("Hook", "applyListToTask", true).Panic("cannot find close list")
		return nil, errors.New("cannot find close list")
	}

	noneList = parentProject.NoneList()
	if noneList == nil {
		logging.SharedInstance().MethodInfo("Hook", "applyListToTask", true).Panic("cannot find none list")
		return nil, errors.New("cannot find none list")
	}

	githubLabels := GithubLabels(issue, parentProject.Lists())

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
