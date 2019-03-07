package repositories

import "database/sql"

// DummyTask defines dummy repository.
type DummyTask struct {
	ListID    int64
	ProjectID int64
	UserID    int64
}

// Find returns a dummy task.
func (d *DummyTask) Find(id int64) (int64, int64, int64, int64, sql.NullInt64, string, string, bool, sql.NullString, error) {
	return id, d.ListID, d.ProjectID, d.UserID, sql.NullInt64{Int64: 1, Valid: true}, "title", "description", false, sql.NullString{String: "", Valid: false}, nil
}

// FindByIssueNumber returns a dummy task.
func (d *DummyTask) FindByIssueNumber(projectID int64, issueNumber int) (int64, int64, int64, int64, sql.NullInt64, string, string, bool, sql.NullString, error) {
	return 1, d.ListID, projectID, d.UserID, sql.NullInt64{Int64: int64(issueNumber), Valid: true}, "title", "description", false, sql.NullString{String: "https://github.com/h3poteto/fascia/issues/1", Valid: true}, nil
}

// Create returns taskID.
func (d *DummyTask) Create(listID, projectID, userID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString) (int64, error) {
	return 1, nil
}

// Update returns no error.
func (d *DummyTask) Update(id, listID, projectID, userID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString) error {
	return nil
}

// ChangeList returns no error.
func (d *DummyTask) ChangeList(id, lisID int64, prevToTaskID *int64) error {
	return nil
}

// Delete returns no error.
func (d *DummyTask) Delete(id int64) error {
	return nil
}

// Tasks returns dummy tasks.
func (d *DummyTask) Tasks(listID int64) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	l := map[string]interface{}{
		"id":          1,
		"listID":      listID,
		"projectID":   d.ProjectID,
		"userID":      d.UserID,
		"issueNumber": sql.NullInt64{},
		"title":       "title",
		"description": "description",
		"pullRequest": false,
		"htmlURL":     sql.NullString{},
	}
	result = append(result, l)
	return result, nil
}

// NonIssueTasks returns dummy tasks.
func (d *DummyTask) NonIssueTasks(projectID, userID int64) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	l := map[string]interface{}{
		"id":          2,
		"listID":      d.ListID,
		"projectID":   projectID,
		"userID":      userID,
		"issueNumber": sql.NullInt64{},
		"title":       "title",
		"description": "description",
		"pullRequest": false,
		"htmlURL":     sql.NullString{},
	}
	result = append(result, l)
	return result, nil
}
