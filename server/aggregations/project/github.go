package project

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/hub"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/models/list"
	"github.com/h3poteto/fascia/server/models/task"

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

func (p *Project) labelUpdate(l *list.ListStruct, labels []*github.Label) error {
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
func (p *Project) TaskLoadFromGithub(issues []*github.Issue) error {
	for _, issue := range issues {
		targetTask, _ := task.FindByIssueNumber(p.ProjectModel.ID, *issue.Number)

		err := p.TaskApplyLabel(targetTask, issue)
		if err != nil {
			return err
		}
	}
	return nil
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

func (p *Project) ReacquireIssue(issue *github.Issue) (*github.Issue, error) {
	if len(issue.Labels) > 0 {
		return issue, nil
	}
	oauthToken, err := p.OauthToken()
	if err != nil {
		return nil, err
	}

	repo, err := p.Repository()
	if err != nil {
		return nil, err
	}
	return hub.GetGithubIssue(oauthToken, repo, *issue.Number)
}

func (p *Project) TaskApplyLabel(targetTask *task.TaskStruct, issue *github.Issue) error {
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
	if err := issueTask.Update(nil, nil); err != nil {
		return err
	}
	return nil
}

func (p *Project) ReopenTask(targetTask *task.TaskStruct, issue *github.Issue) error {
	issueTask, err := p.applyListToTask(targetTask, issue)
	if err != nil {
		return err
	}
	if err := issueTask.Update(nil, nil); err != nil {
		return err
	}
	return nil
}

func (p *Project) CreateNewTask(issue *github.Issue) error {
	issueTask := task.NewTask(
		0,
		0,
		p.ProjectModel.ID,
		p.ProjectModel.UserID,
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
	if err := issueTask.Save(nil, nil); err != nil {
		return err
	}
	return nil

}

func (p *Project) applyListToTask(issueTask *task.TaskStruct, issue *github.Issue) (*task.TaskStruct, error) {
	// close noneの用意
	var closedList, noneList *list.ListStruct
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

func (u *Project) applyIssueInfoToTask(targetTask *task.TaskStruct, issue *github.Issue) (*task.TaskStruct, error) {
	if targetTask == nil {
		return nil, errors.New("target task is required")
	}
	targetTask.Title = *issue.Title
	targetTask.Description = *issue.Body
	targetTask.HTMLURL = sql.NullString{String: *issue.HTMLURL, Valid: true}
	return targetTask, nil
}
