package board

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/hub"
	"github.com/h3poteto/fascia/lib/modules/logging"
	rediscli "github.com/h3poteto/fascia/lib/modules/redis"
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/domains/task"
	"github.com/pkg/errors"
)

// applyIssueChanges apply issue changes to task
func applyIssueChanges(p *project.Project, body github.IssuesEvent) error {
	logging.SharedInstance().MethodInfo("Project", "applyIssueChanges").Debugf("project: %+v", p)
	logging.SharedInstance().MethodInfo("Project", "applyIssueChanges").Debugf("issues event: %+v", *body.Issue)
	infra := InjectTaskRepository()
	// When the task does not exist, we create a new task.
	// So if this method returns an error, we can ignore it.
	targetTask, _ := infra.FindByIssueNumber(p.ID, *body.Issue.Number)

	// create時点ではlabelsが空の状態でhookが飛んできている場合がある
	// editedの場合であってもwebhookにはchangeだけしか載っておらず，最新の状態は載っていない場合がある
	// そのため一度issueの情報を取得し直す必要がある
	issue, err := reacquireIssue(p, body.Issue)
	if err != nil {
		return err
	}
	switch *body.Action {
	case "opened", "reopened":
		if targetTask == nil {
			err = createNewTask(p, issue)
		} else {
			err = reopenTask(p, targetTask, issue)
		}
	case "closed", "labeled", "unlabeled", "edited":
		err = taskApplyLabel(p, targetTask, issue)
	}
	return err
}

// reacquireIssue get again a issue from github
func reacquireIssue(p *project.Project, issue *github.Issue) (*github.Issue, error) {
	repo := InjectProjectRepository()
	oauthToken, err := repo.OauthToken(p.ID)
	if err != nil {
		return nil, err
	}

	r, err := ProjectRepository(p)
	if err != nil {
		return nil, errors.Wrap(err, "can not find repository")
	}
	return r.GetGithubIssue(oauthToken, *issue.Number)
}

// createNewTask create a task from github issue
func createNewTask(p *project.Project, issue *github.Issue) error {
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

	issueTask, err := applyListToTask(p, issueTask, issue)
	if err != nil {
		return err
	}
	infra := InjectTaskRepository()
	if _, err := infra.Create(issueTask.ListID, issueTask.ProjectID, issueTask.UserID, issueTask.IssueNumber, issueTask.Title, issueTask.Description, issueTask.PullRequest, issueTask.HTMLURL); err != nil {
		return err
	}
	return nil
}

// reopenTask open a task according to github issue
// It is not necessary to change a task status in Database
// It is enough to update issue information, and change label
func reopenTask(p *project.Project, targetTask *task.Task, issue *github.Issue) error {
	issueTask, err := applyListToTask(p, targetTask, issue)
	if err != nil {
		return err
	}
	infra := InjectTaskRepository()
	err = infra.Update(
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

// FetchGithub fetch all lists and all tasks
func FetchGithub(p *project.Project) (bool, error) {
	repo, err := ProjectRepository(p)
	if err != nil {
		return false, err
	}

	infra := InjectProjectRepository()
	oauthToken, err := infra.OauthToken(p.ID)
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
	err = listLoadFromGithub(p, labels)
	if err != nil {
		return false, err
	}

	// Import tasks from issues.
	openIssues, closedIssues, err := repo.GetGithubIssues(oauthToken)
	if err != nil {
		return false, err
	}

	err = taskLoadFromGithub(p, append(openIssues, closedIssues...))
	if err != nil {
		return false, err
	}

	//-------------------------------------
	// Export lists and tasks to Github.
	//------------------------------------
	taskInfra := InjectTaskRepository()
	tasks, err := taskInfra.NonIssueTasks(p.ID, p.UserID)
	if err != nil {
		return false, err
	}
	for _, t := range tasks {
		listRepo := InjectListRepository()
		l, err := listRepo.FindByTaskID(t.ID)
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

// applyPullRequestChanges apply issue changes to task
func applyPullRequestChanges(p *project.Project, body github.PullRequestEvent) error {
	// When the task does not exist, we create a new task.
	// So if this method returns an error, we can ignore it.
	taskInfra := InjectTaskRepository()
	targetTask, _ := taskInfra.FindByIssueNumber(p.ID, *body.Number)

	infra := InjectProjectRepository()
	oauthToken, err := infra.OauthToken(p.ID)
	if err != nil {
		return err
	}

	repo, err := ProjectRepository(p)
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
			err = createNewTask(p, issue)
		} else {
			err = reopenTask(p, targetTask, issue)
		}
	case "closed", "labeled", "unlabeled", "edited":
		err = taskApplyLabel(p, targetTask, issue)
	}
	return err
}

// taskApplyLabel change list to a task according to github labels
func taskApplyLabel(p *project.Project, targetTask *task.Task, issue *github.Issue) error {
	if targetTask == nil {
		err := createNewTask(p, issue)
		if err != nil {
			return err
		}
		return nil
	}
	issueTask, err := applyListToTask(p, targetTask, issue)
	if err != nil {
		return err
	}
	issueTask, err = applyIssueInfoToTask(p, issueTask, issue)
	if err != nil {
		return err
	}
	infra := InjectTaskRepository()
	err = infra.Update(
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

func applyListToTask(p *project.Project, issueTask *task.Task, issue *github.Issue) (*task.Task, error) {
	// close noneの用意
	var closedList, noneList *list.List
	listRepo := InjectListRepository()
	lists, err := listRepo.Lists(p.ID)
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

	noneList, err = listRepo.NoneList(p.ID)
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

// taskLoadFromGithub load tasks from github issues
func taskLoadFromGithub(p *project.Project, issues []*github.Issue) error {
	infra := InjectTaskRepository()
	for _, issue := range issues {
		targetTask, _ := infra.FindByIssueNumber(p.ID, *issue.Number)

		err := taskApplyLabel(p, targetTask, issue)
		if err != nil {
			return err
		}
	}
	return nil
}

// listLoadFromGithub load lists from github labels
func listLoadFromGithub(p *project.Project, labels []*github.Label) error {
	listRepo := InjectListRepository()
	lists, err := listRepo.Lists(p.ID)
	if err != nil {
		return err
	}
	for _, l := range lists {
		if err := labelUpdate(l, labels); err != nil {
			return err
		}
	}
	return nil
}

func labelUpdate(l *list.List, labels []*github.Label) error {
	for _, label := range labels {
		if strings.ToLower(*label.Name) == strings.ToLower(l.Title.String) {
			title := sql.NullString{String: l.Title.String, Valid: true}
			color := sql.NullString{String: *label.Color, Valid: true}
			if err := l.Update(title, color, l.Option); err != nil {
				return err
			}
			repo := InjectListRepository()
			if err := repo.Update(l); err != nil {
				return err
			}
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

// InjectRedis returns a redis client instance.
func InjectRedis() *redis.Client {
	return rediscli.SharedInstance().Client
}

// GetAllRepositories gets all repositoreis from github related the oauth token.
func GetAllRepositories(oauthToken string) ([]*github.Repository, error) {
	cli := InjectRedis()
	val, err := cli.LRange(oauthToken, 0, -1).Result()
	logging.SharedInstance().MethodInfo("board", "GetAllRepositories").Debugf("redis values: %+v", val)
	if err == nil && len(val) > 0 {
		var res []*github.Repository
		for _, jsonStr := range val {
			jsonBytes := ([]byte)(jsonStr)
			var r github.Repository
			if err := json.Unmarshal(jsonBytes, &r); err != nil {
				return nil, errors.Wrap(err, "Unmarshal error")
			}
			res = append(res, &r)
		}
		return res, nil
	}
	repositories, err := hub.New(oauthToken).AllRepositories()
	if err != nil {
		return nil, err
	}
	go func() {
		for _, repository := range repositories {
			jsonBytes, err := json.Marshal(repository)
			if err != nil {
				logging.SharedInstance().MethodInfo("board", "GetAllRepositories").Error(err)
				return
			}
			err = cli.RPush(oauthToken, string(jsonBytes)).Err()
			if err != nil {
				logging.SharedInstance().MethodInfo("board", "GetAllRepositories").Error(err)
				return
			}
		}
		// TODO: Extend the expire after you implement refresh function.
		err := cli.Expire(oauthToken, 48*time.Hour).Err()
		if err != nil {
			logging.SharedInstance().MethodInfo("board", "GetAllRepositories").Error(err)
			return
		}
		return
	}()
	return repositories, nil
}
