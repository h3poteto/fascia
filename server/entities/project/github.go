package project

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/hub"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/entities/list"
	"github.com/h3poteto/fascia/server/entities/task"

	"database/sql"
	"strings"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

// ListLoadFromGithub load lists from github labels
func (p *Project) ListLoadFromGithub(labels []*github.Label) error {
	lists, err := p.Lists()
	if err != nil {
		return err
	}
	for _, l := range lists {
		if err := p.labelUpdate(l, labels); err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) labelUpdate(l *list.List, labels []*github.Label) error {
	for _, label := range labels {
		if strings.ToLower(*label.Name) == strings.ToLower(l.Title.String) {
			if err := l.Update(l.Title.String, *label.Color, l.ListOptionID.Int64); err != nil {
				return err
			}
		}
	}
	return nil
}

// TaskLoadFromGithub load tasks from github issues
func (p *Project) TaskLoadFromGithub(issues []*github.Issue) error {
	for _, issue := range issues {
		targetTask, _ := task.FindByIssueNumber(p.ID, *issue.Number)

		err := p.TaskApplyLabel(targetTask, issue)
		if err != nil {
			return err
		}
	}
	return nil
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

// ReacquireIssue get again a issue from github
func (p *Project) ReacquireIssue(issue *github.Issue) (*github.Issue, error) {
	oauthToken, err := p.OauthToken()
	if err != nil {
		return nil, err
	}

	repo, err := p.Repository()
	if err != nil {
		return nil, errors.New("can not find repository")
	}
	return repo.GetGithubIssue(oauthToken, *issue.Number)
}

// TaskApplyLabel change list to a task according to github labels
func (p *Project) TaskApplyLabel(targetTask *task.Task, issue *github.Issue) error {
	if targetTask == nil {
		err := p.CreateNewTask(issue)
		if err != nil {
			return err
		}
		return nil
	}
	issueTask, err := p.applyListToTask(targetTask, issue)
	if err != nil {
		return err
	}
	issueTask, err = p.applyIssueInfoToTask(issueTask, issue)
	if err != nil {
		return err
	}
	err = issueTask.Update(
		issueTask.ListID,
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

// ReopenTask open a task according to github issue
// It is not necessary to change a task status in Database
// It is enough to update issue information, and change label
func (p *Project) ReopenTask(targetTask *task.Task, issue *github.Issue) error {
	issueTask, err := p.applyListToTask(targetTask, issue)
	if err != nil {
		return err
	}
	err = issueTask.Update(
		issueTask.ListID,
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

// CreateNewTask create a task from github issue
func (p *Project) CreateNewTask(issue *github.Issue) error {
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

	issueTask, err := p.applyListToTask(issueTask, issue)
	if err != nil {
		return err
	}
	if err := issueTask.Save(); err != nil {
		return err
	}
	return nil
}

func (p *Project) applyListToTask(issueTask *task.Task, issue *github.Issue) (*task.Task, error) {
	// close noneの用意
	var closedList, noneList *list.List
	lists, err := p.Lists()
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

	noneList, err = p.NoneList()
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

func (p *Project) applyIssueInfoToTask(targetTask *task.Task, issue *github.Issue) (*task.Task, error) {
	if targetTask == nil {
		return nil, errors.New("target task is required")
	}
	targetTask.Title = *issue.Title
	targetTask.Description = *issue.Body
	targetTask.HTMLURL = sql.NullString{String: *issue.HTMLURL, Valid: true}
	return targetTask, nil
}
