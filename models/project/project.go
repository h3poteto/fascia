package project

import (
	"../../config"
	"../../modules/hub"
	"../../modules/logging"
	"../db"
	"../list"
	"../list_option"
	"../repository"
	"../task"
	"database/sql"
	"errors"
)

type Project interface {
	Lists() []*list.ListStruct
	Save() bool
}

type ProjectStruct struct {
	Id           int64
	UserId       int64
	Title        string
	Description  string
	RepositoryId sql.NullInt64
	database     db.DB
}

func NewProject(id int64, userID int64, title string, description string, repositoryId sql.NullInt64) *ProjectStruct {
	if userID == 0 {
		return nil
	}

	project := &ProjectStruct{Id: id, UserId: userID, Title: title, Description: description, RepositoryId: repositoryId}
	project.Initialize()
	return project
}

func FindProject(projectID int64) *ProjectStruct {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var id, userID int64
	var repositoryID sql.NullInt64
	var title string
	var description string
	err := table.QueryRow("select id, user_id, repository_id, title, description from projects where id = ?;", projectID).Scan(&id, &userID, &repositoryID, &title, &description)
	if err != nil {
		logging.SharedInstance().MethodInfo("Project", "FindProject", true).Errorf("cannot find project: %v", err)
		return nil
	}
	project := NewProject(id, userID, title, description, repositoryID)
	return project
}

func Create(userID int64, title string, description string, repositoryID int64, repositoryOwner string, repositoryName string, oauthToken sql.NullString) (p *ProjectStruct, e error) {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()
	tx, _ := table.Begin()
	defer func() {
		if err := recover(); err != nil {
			logging.SharedInstance().MethodInfo("Project", "Create", true).Error("unexpected error")
			tx.Rollback()
			e = errors.New("unexpected error")
			p = nil
		}
	}()

	var repoID sql.NullInt64
	var repo *repository.RepositoryStruct
	if repositoryID != 0 {
		repo = repository.NewRepository(0, repositoryID, repositoryOwner, repositoryName)
		if !repo.Save() {
			tx.Rollback()
			logging.SharedInstance().MethodInfo("Project", "Create", true).Error("failed to save repository")
			return nil, errors.New("repository save error")
		}
		repoID = sql.NullInt64{Int64: repo.Id, Valid: true}
	}

	project := NewProject(0, userID, title, description, repoID)
	if !project.Save() {
		tx.Rollback()
		logging.SharedInstance().MethodInfo("Project", "Create", true).Error("failed to save project")
		return nil, errors.New("failed to save project")
	}

	// 初期リストの準備
	closeListOption := list_option.FindByAction("close")
	if closeListOption == nil {
		tx.Rollback()
		logging.SharedInstance().MethodInfo("Project", "Create", true).Error("cannot find close list option")
		return nil, errors.New("failed to find close list option")
	}
	todo := list.NewList(0, project.Id, userID, config.Element("init_list").(map[interface{}]interface{})["todo"].(string), "f37b1d", sql.NullInt64{})
	inprogress := list.NewList(0, project.Id, userID, config.Element("init_list").(map[interface{}]interface{})["inprogress"].(string), "5eb95e", sql.NullInt64{})
	done := list.NewList(0, project.Id, userID, config.Element("init_list").(map[interface{}]interface{})["done"].(string), "333333", sql.NullInt64{Int64: closeListOption.Id, Valid: true})
	none := list.NewList(0, project.Id, userID, config.Element("init_list").(map[interface{}]interface{})["none"].(string), "ffffff", sql.NullInt64{})
	if !none.Save(nil, nil) {
		tx.Rollback()
		logging.SharedInstance().MethodInfo("Project", "Create", true).Error("failed to save none list")
		return nil, errors.New("failed to save none list")
	}

	if project.RepositoryId.Valid {
		if !todo.Save(repo, &oauthToken) {
			tx.Rollback()
			logging.SharedInstance().MethodInfo("Project", "Create", true).Error("failed to save todo list")
			return nil, errors.New("failed to save todo list")
		}
		if !inprogress.Save(repo, &oauthToken) {
			tx.Rollback()
			logging.SharedInstance().MethodInfo("Project", "Create", true).Error("failed to save inprogress list")
			return nil, errors.New("failed to save inprogress list")
		}
		if !done.Save(repo, &oauthToken) {
			tx.Rollback()
			logging.SharedInstance().MethodInfo("Project", "Create", true).Error("failed to save done list")
			return nil, errors.New("failed to save done list")
		}
	} else {
		if !todo.Save(nil, nil) {
			tx.Rollback()
			logging.SharedInstance().MethodInfo("Project", "Create", true).Error("failed to save todo list")
			return nil, errors.New("failed to save todo list")
		}
		if !inprogress.Save(nil, nil) {
			tx.Rollback()
			logging.SharedInstance().MethodInfo("Project", "Create", true).Error("failed to save inprogress list")
			return nil, errors.New("failed to save inprogress list")
		}
		if !done.Save(nil, nil) {
			tx.Rollback()
			logging.SharedInstance().MethodInfo("Project", "Create", true).Error("failed to save done list")
			return nil, errors.New("failed to save done list")
		}
	}
	tx.Commit()
	return project, nil
}

func (u *ProjectStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

func (u *ProjectStruct) Save() bool {
	table := u.database.Init()
	defer table.Close()

	result, err := table.Exec("insert into projects (user_id, repository_id, title, description, created_at) values (?, ?, ?, ?, now());", u.UserId, u.RepositoryId, u.Title, u.Description)
	if err != nil {
		logging.SharedInstance().MethodInfo("Project", "Save", true).Errorf("failed to save project: %v", err)
		return false
	}
	u.Id, _ = result.LastInsertId()
	return true
}

func (u *ProjectStruct) Update(title string, description string) bool {
	table := u.database.Init()
	defer table.Close()

	u.Title = title
	u.Description = description
	_, err := table.Exec("update projects set title = ?, description = ? where id = ?;", u.Title, u.Description, u.Id)
	if err != nil {
		logging.SharedInstance().MethodInfo("Project", "Update", true).Errorf("failed to update project: %v", err)
		return false
	}

	return true
}

func (u *ProjectStruct) Lists() []*list.ListStruct {
	table := u.database.Init()
	defer table.Close()

	rows, _ := table.Query("select id, project_id, user_id, title, color, list_option_id from lists where project_id = ? and title != ?;", u.Id, config.Element("init_list").(map[interface{}]interface{})["none"].(string))
	var slice []*list.ListStruct
	for rows.Next() {
		var id, projectID, userID int64
		var title, color sql.NullString
		var optionID sql.NullInt64
		err := rows.Scan(&id, &projectID, &userID, &title, &color, &optionID)
		if err != nil {
			logging.SharedInstance().MethodInfo("Project", "Lists").Panic(err)
		}
		if projectID == u.Id && title.Valid {
			l := list.NewList(id, projectID, userID, title.String, color.String, optionID)
			slice = append(slice, l)
		}
	}
	return slice
}

func (u *ProjectStruct) NoneList() *list.ListStruct {
	table := u.database.Init()
	defer table.Close()

	var id, projectID, userID int64
	var title, color sql.NullString
	var optionID sql.NullInt64
	err := table.QueryRow("select id, project_id, user_id, title, color, list_option_id from lists where project_id = ? and title = ?;", u.Id, config.Element("init_list").(map[interface{}]interface{})["none"].(string)).Scan(&id, &projectID, &userID, &title, &color, &optionID)
	if err != nil {
		// noneが存在しないということはProjectsController#Createがうまく行ってないので，そっちでエラーハンドリングしてほしい
		logging.SharedInstance().MethodInfo("Project", "NoneList").Panic(err)
	}
	if projectID == u.Id && title.Valid {
		return list.NewList(id, projectID, userID, title.String, color.String, optionID)
	}
	return nil
}

func (u *ProjectStruct) Repository() *repository.RepositoryStruct {
	table := u.database.Init()
	defer table.Close()

	var id, repositoryId int64
	var owner, name sql.NullString
	err := table.QueryRow("select repositories.id, repositories.repository_id, repositories.owner, repositories.name from projects inner join repositories on repositories.id = projects.repository_id where projects.id = ?;", u.Id).Scan(&id, &repositoryId, &owner, &name)
	if err != nil {
		logging.SharedInstance().MethodInfo("project", "Repository").Infof("cannot find repository: %v", err)
		return nil
	}
	if id == u.RepositoryId.Int64 && owner.Valid {
		r := repository.NewRepository(id, repositoryId, owner.String, name.String)
		return r
	} else {
		logging.SharedInstance().MethodInfo("project", "Repository", true).Error("repository owner discord from project owner")
		return nil
	}
}

func (u *ProjectStruct) FetchGithub() (bool, error) {
	table := u.database.Init()
	defer table.Close()

	var oauthToken sql.NullString
	err := table.QueryRow("select users.oauth_token from projects left join users on users.id = projects.user_id where projects.id = ?;", u.Id).Scan(&oauthToken)
	if err != nil {
		logging.SharedInstance().MethodInfo("project", "FetchGithub", true).Errorf("oauth_token select error: %v", err)
		return false, err
	}
	if !oauthToken.Valid {
		logging.SharedInstance().MethodInfo("project", "FetchGithub").Info("oauth token is not valid")
		return false, errors.New("oauth token is required")
	}
	repo := u.Repository()
	// user自体はgithub連携していても，projectが連携していない可能性もあるのでチェック
	if repo == nil {
		logging.SharedInstance().MethodInfo("project", "FetchGithub").Info("repository related project is nil")
		return false, errors.New("project did not related to repository")
	}

	openIssues, closedIssues, err := hub.GetGithubIssues(oauthToken.String, repo)
	if err != nil {
		return false, err
	}
	var closedList *list.ListStruct
	for _, list := range u.Lists() {
		// closeのリストは用意しておく
		if list.Title.Valid && list.Title.String == config.Element("init_list").(map[interface{}]interface{})["done"].(string) {
			closedList = list
		}
	}
	if closedList == nil {
		logging.SharedInstance().MethodInfo("Project", "FetchGithub").Panic("cannot find close list")
		return false, errors.New("cannot find close list")
	}
	noneList := u.NoneList()
	if noneList == nil {
		logging.SharedInstance().MethodInfo("Project", "FetchGithub").Panic("cannot find none list")
		return false, errors.New("cannot find none list")
	}

	for _, issue := range append(openIssues, closedIssues...) {
		var githubLabels []list.ListStruct
		for _, label := range issue.Labels {
			for _, list := range u.Lists() {
				// 紐付いているlabelのlistを持っている時
				if list.Title.Valid && list.Title.String == *label.Name {
					githubLabels = append(githubLabels, *list)
				}
			}
		}
		issueTask, err := task.FindByIssueNumber(*issue.Number)
		if err != nil && issueTask == nil {
			issueTask = task.NewTask(0, 0, u.UserId, sql.NullInt64{Int64: int64(*issue.Number), Valid: true}, *issue.Title, *issue.Body)
		}
		if len(githubLabels) == 1 {
			// 一つのlistだけが該当するとき
			issueTask.ListId = githubLabels[0].Id
		} else if len(githubLabels) > 1 {
			// 複数のlistが該当するとき
			issueTask.ListId = githubLabels[0].Id
		} else {
			// ついているlabelのlistを持ってない時
			if *issue.State == "open" {
				issueTask.ListId = noneList.Id
			} else {
				issueTask.ListId = closedList.Id
			}
		}
		// ここはgithub側への同期不要
		if issueTask.Id == 0 {
			if !issueTask.Save(nil, nil) {
				logging.SharedInstance().MethodInfo("Project", "FetchGithub", true).Error("failed to save task")
				return false, errors.New("failed to save task")
			}
		} else {
			issueTask.Title = *issue.Title
			issueTask.Description = *issue.Body
			if !issueTask.Update(nil, nil) {
				logging.SharedInstance().MethodInfo("Project", "FetchGithub", true).Error("failed to update task")
				return false, errors.New("failed to update task")
			}
		}
	}
	// github側へ同期
	rows, err := table.Query("select tasks.title, tasks.description, lists.title, lists.color from tasks left join lists on lists.id = tasks.list_id where tasks.user_id = ? and tasks.issue_number IS NULL;", u.UserId)
	if err != nil {
		logging.SharedInstance().MethodInfo("Project", "FetchGithub").Panic(err)
		return false, err
	}
	for rows.Next() {
		var title, description string
		var listTitle, listColor sql.NullString
		err := rows.Scan(&title, &description, &listTitle, &listColor)
		if err != nil {
			return false, err
		}
		label, err := hub.CheckLabelPresent(oauthToken.String, repo, &listTitle.String)
		if err != nil {
			return false, err
		}
		if label == nil {
			label, err = hub.CreateGithubLabel(oauthToken.String, repo, &listTitle.String, &listColor.String)
			if err != nil {
				return false, errors.New("cannot create github label")
			}
		}
		// ここcreateだけでなくupdateも考慮したほうが良いのではと思ったが，そもそも現状fasciaにはtaskのupdateアクションがないので，updateされることはありえない．そのため，未実装でも問題はない．
		// todo: task#update実装時にはここも実装すること
		_, err = hub.CreateGithubIssue(oauthToken.String, repo, []string{*label.Name}, &title, &description)
		if err != nil {
			return false, errors.New("cannot create github issue")
		}
	}

	return true, nil
}
