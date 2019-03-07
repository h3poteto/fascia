package repositories

import "database/sql"

// DummyProject defines dummy repository.
type DummyProject struct {
	UserID int64
}

// Find return a dummy project.
func (d *DummyProject) Find(id int64) (int64, int64, string, string, sql.NullInt64, bool, bool, error) {
	return id, d.UserID, "title", "description", sql.NullInt64{}, true, true, nil
}

// FindByRepositoryID returns dummy projects.
func (d *DummyProject) FindByRepositoryID(repoID int64) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	p := map[string]interface{}{
		"id":               1,
		"userID":           d.UserID,
		"title":            "title",
		"description":      "description",
		"repositoryID":     sql.NullInt64{Int64: repoID, Valid: true},
		"showIssues":       true,
		"showPullRequests": true,
	}
	result = append(result, p)
	return result, nil
}

// Create returns a projectID.
func (d *DummyProject) Create(userID int64, title, description string, repoID sql.NullInt64, showIssues, showPullRequests bool, tx *sql.Tx) (int64, error) {
	return 1, nil
}

// Update returns no error.
func (d *DummyProject) Update(id, userID int64, title, description string, repoID sql.NullInt64, showIssues, showPullRequests bool) error {
	return nil
}

// Delete returns no error.
func (d *DummyProject) Delete(id int64) error {
	return nil
}

// Projects returns dummy projects.
func (d *DummyProject) Projects(userID int64) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	p := map[string]interface{}{
		"id":               1,
		"userID":           userID,
		"title":            "title",
		"description":      "description",
		"repositoryID":     sql.NullInt64{},
		"showIssues":       true,
		"showPullRequests": true,
	}
	result = append(result, p)
	return result, nil
}

// OauthToken returns a dummy token.
func (d *DummyProject) OauthToken(id int64) (string, error) {
	return "dummy", nil
}
