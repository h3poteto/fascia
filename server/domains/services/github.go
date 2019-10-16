package services

import (
	"database/sql"
	"strings"

	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/hub"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/domains/repo"
	"github.com/h3poteto/fascia/server/domains/task"
	"github.com/pkg/errors"
)

// FetchGithub fetch all lists and all tasks
func FetchGithub(p *project.Project, projectInfra project.Repository, listInfra list.Repository, taskInfra task.Repository, repoInfra repo.Repository) (bool, error) {
	repo, err := repoInfra.FindByProjectID(p.ID)
	if err != nil {
		return false, err
	}

	oauthToken, err := projectInfra.OauthToken(p.ID)
	if err != nil {
		return false, err
	}

	//------------------------------------
	// Import lists and tasks from Github.
	//------------------------------------
	// Import lists from labels.
	labels, err := repo.ListLabels(oauthToken)
	if err != nil {
		return false, err
	}
	err = listLoadFromGithub(p, labels, listInfra)
	if err != nil {
		return false, err
	}

	// Import tasks from issues.
	openIssues, closedIssues, err := repo.GetGithubIssues(oauthToken)
	if err != nil {
		return false, err
	}

	err = taskLoadFromGithub(p, append(openIssues, closedIssues...), listInfra, taskInfra)
	if err != nil {
		return false, err
	}

	//-------------------------------------
	// Export lists and tasks to Github.
	//------------------------------------
	tasks, err := taskInfra.NonIssueTasks(p.ID, p.UserID)
	if err != nil {
		return false, err
	}
	for _, t := range tasks {
		l, err := listInfra.FindByTaskID(t.ID)
		if err != nil {
			return false, err
		}
		label, err := repo.CheckLabelPresent(oauthToken, l.Title.String)
		if err != nil {
			return false, err
		}
		if label == nil {
			label, err = repo.CreateGithubLabel(oauthToken, l.Title.String, l.Color.String)
			if err != nil {
				return false, err
			}
		}

		_, err = repo.CreateGithubIssue(oauthToken, t.Title, t.Description, []string{*label.Name})
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// listLoadFromGithub load lists from github labels
func listLoadFromGithub(p *project.Project, labels []*github.Label, listInfra list.Repository) error {
	lists, err := listInfra.Lists(p.ID)
	if err != nil {
		return err
	}
	for _, l := range lists {
		if err := labelUpdate(l, labels, listInfra); err != nil {
			return err
		}
	}
	return nil
}

// taskLoadFromGithub load tasks from github issues
func taskLoadFromGithub(p *project.Project, issues []*github.Issue, listInfra list.Repository, taskInfra task.Repository) error {
	for _, issue := range issues {
		targetTask, _ := taskInfra.FindByIssueNumber(p.ID, *issue.Number)

		err := taskApplyLabel(p, targetTask, issue, listInfra, taskInfra)
		if err != nil {
			return err
		}
	}
	return nil
}

func labelUpdate(l *list.List, labels []*github.Label, listInfra list.Repository) error {
	for _, label := range labels {
		if strings.ToLower(*label.Name) == strings.ToLower(l.Title.String) {
			title := sql.NullString{String: l.Title.String, Valid: true}
			color := sql.NullString{String: *label.Color, Valid: true}
			if err := l.Update(title, color, l.Option); err != nil {
				return err
			}
			if err := listInfra.Update(l); err != nil {
				return err
			}
		}
	}
	return nil
}

// taskApplyLabel change list to a task according to github labels
func taskApplyLabel(p *project.Project, targetTask *task.Task, issue *github.Issue, listInfra list.Repository, taskInfra task.Repository) error {
	if targetTask == nil {
		err := createNewTask(p, issue, listInfra, taskInfra)
		if err != nil {
			return err
		}
		return nil
	}
	issueTask, err := applyListToTask(p, targetTask, issue, listInfra)
	if err != nil {
		return err
	}
	issueTask, err = applyIssueInfoToTask(p, issueTask, issue)
	if err != nil {
		return err
	}
	err = taskInfra.Update(
		issueTask.ID,
		issueTask.ListID,
		issueTask.ProjectID,
		issueTask.UserID,
		issueTask.IssueNumber,
		issueTask.Title,
		issueTask.Description,
		issueTask.PullRequest,
		issueTask.HTMLURL,
	)
	if err != nil {
		return err
	}
	return nil
}

// createNewTask create a task from github issue
func createNewTask(p *project.Project, issue *github.Issue, listInfra list.Repository, taskInfra task.Repository) error {
	issueTask := task.New(
		0,
		0,
		p.ID,
		p.UserID,
		sql.NullInt64{Int64: int64(*issue.Number), Valid: true},
		*issue.Title,
		*issue.Body,
		hub.IsPullRequest(issue),
		sql.NullString{String: *issue.HTMLURL, Valid: true},
	)

	issueTask, err := applyListToTask(p, issueTask, issue, listInfra)
	if err != nil {
		return err
	}
	if _, err := taskInfra.Create(issueTask.ListID, issueTask.ProjectID, issueTask.UserID, issueTask.IssueNumber, issueTask.Title, issueTask.Description, issueTask.PullRequest, issueTask.HTMLURL); err != nil {
		return err
	}
	return nil
}

func applyListToTask(p *project.Project, issueTask *task.Task, issue *github.Issue, listInfra list.Repository) (*task.Task, error) {
	// close noneの用意
	var closedList, noneList *list.List
	lists, err := listInfra.Lists(p.ID)
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

	noneList, err = listInfra.NoneList(p.ID)
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

func applyIssueInfoToTask(p *project.Project, targetTask *task.Task, issue *github.Issue) (*task.Task, error) {
	if targetTask == nil {
		return nil, errors.New("target task is required")
	}
	targetTask.Title = *issue.Title
	targetTask.Description = *issue.Body
	targetTask.HTMLURL = sql.NullString{String: *issue.HTMLURL, Valid: true}
	return targetTask, nil
}

// applyIssueChanges apply issue changes to task
func applyIssueChanges(p *project.Project, body github.IssuesEvent, projectInfra project.Repository, listInfra list.Repository, taskInfra task.Repository, repoInfra repo.Repository) error {
	logging.SharedInstance().MethodInfo("Project", "applyIssueChanges").Debugf("project: %+v", p)
	logging.SharedInstance().MethodInfo("Project", "applyIssueChanges").Debugf("issues event: %+v", *body.Issue)
	// When the task does not exist, we create a new task.
	// So if this method returns an error, we can ignore it.
	targetTask, _ := taskInfra.FindByIssueNumber(p.ID, *body.Issue.Number)

	// create時点ではlabelsが空の状態でhookが飛んできている場合がある
	// editedの場合であってもwebhookにはchangeだけしか載っておらず，最新の状態は載っていない場合がある
	// そのため一度issueの情報を取得し直す必要がある
	issue, err := reacquireIssue(p, body.Issue, projectInfra, repoInfra)
	if err != nil {
		return err
	}
	switch *body.Action {
	case "opened", "reopened":
		if targetTask == nil {
			err = createNewTask(p, issue, listInfra, taskInfra)
		} else {
			err = reopenTask(p, targetTask, issue, listInfra, taskInfra)
		}
	case "closed", "labeled", "unlabeled", "edited":
		err = taskApplyLabel(p, targetTask, issue, listInfra, taskInfra)
	}
	return err
}

// reacquireIssue get again a issue from github
func reacquireIssue(p *project.Project, issue *github.Issue, projectInfra project.Repository, repoInfra repo.Repository) (*github.Issue, error) {
	oauthToken, err := projectInfra.OauthToken(p.ID)
	if err != nil {
		return nil, err
	}

	r, err := repoInfra.FindByProjectID(p.ID)
	if err != nil {
		return nil, errors.Wrap(err, "can not find repository")
	}
	return r.GetGithubIssue(oauthToken, *issue.Number)
}

// reopenTask open a task according to github issue
// It is not necessary to change a task status in Database
// It is enough to update issue information, and change label
func reopenTask(p *project.Project, targetTask *task.Task, issue *github.Issue, listInfra list.Repository, taskInfra task.Repository) error {
	issueTask, err := applyListToTask(p, targetTask, issue, listInfra)
	if err != nil {
		return err
	}

	err = taskInfra.Update(
		issueTask.ID,
		issueTask.ListID,
		issueTask.ProjectID,
		issueTask.UserID,
		issueTask.IssueNumber,
		issueTask.Title,
		issueTask.Description,
		issueTask.PullRequest,
		issueTask.HTMLURL,
	)
	if err != nil {
		return err
	}
	return nil
}

// applyPullRequestChanges apply issue changes to task
func applyPullRequestChanges(p *project.Project, body github.PullRequestEvent, projectInfra project.Repository, listInfra list.Repository, taskInfra task.Repository, repoInfra repo.Repository) error {
	// When the task does not exist, we create a new task.
	// So if this method returns an error, we can ignore it.
	targetTask, _ := taskInfra.FindByIssueNumber(p.ID, *body.Number)

	oauthToken, err := projectInfra.OauthToken(p.ID)
	if err != nil {
		return err
	}

	repo, err := repoInfra.FindByProjectID(p.ID)
	if err != nil {
		return err
	}
	// createNewTask method requires *github.Issue type, but there is a struct of *github.PullRequest.
	// So we have to get *github.Issue.
	issue, err := repo.GetGithubIssue(oauthToken, *body.Number)
	if err != nil {
		return err
	}

	switch *body.Action {
	case "opened", "reopened":
		if targetTask == nil {
			err = createNewTask(p, issue, listInfra, taskInfra)
		} else {
			err = reopenTask(p, targetTask, issue, listInfra, taskInfra)
		}
	case "closed", "labeled", "unlabeled", "edited":
		err = taskApplyLabel(p, targetTask, issue, listInfra, taskInfra)
	}
	return err
}

// githubLabels get issues which match project's lists from github
func githubLabelLists(issue *github.Issue, projectLists []*list.List) []list.List {
	var githubLabels []list.List
	for _, label := range issue.Labels {
		for _, list := range projectLists {
			if list.Title.Valid && strings.ToLower(list.Title.String) == strings.ToLower(*label.Name) {
				githubLabels = append(githubLabels, *list)
			}
		}
	}
	return githubLabels
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
