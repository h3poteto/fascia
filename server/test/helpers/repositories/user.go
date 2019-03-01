package repositories

import "database/sql"

// DummyUser is a struct for test repository.
type DummyUser struct{}

// Find returns a dummy user.
func (d *DummyUser) Find(id int64) (int64, string, string, sql.NullString, sql.NullString, sql.NullInt64, sql.NullString, sql.NullString, error) {
	return id, "hoge@example.com", "$2a$10$idVJpEyBXEaW0ODQaq2EtecuEm6Mrk.V7YP5lzWmyW11cOg37NIAq", sql.NullString{}, sql.NullString{}, sql.NullInt64{}, sql.NullString{}, sql.NullString{}, nil
}

// FindByEmail returns a dummy user.
func (d *DummyUser) FindByEmail(email string) (int64, string, string, sql.NullString, sql.NullString, sql.NullInt64, sql.NullString, sql.NullString, error) {
	return 1, email, "$2a$10$idVJpEyBXEaW0ODQaq2EtecuEm6Mrk.V7YP5lzWmyW11cOg37NIAq", sql.NullString{}, sql.NullString{}, sql.NullInt64{}, sql.NullString{}, sql.NullString{}, nil
}

// Create returns no error.
func (d *DummyUser) Create(email string, password string, provider sql.NullString, token sql.NullString, uuid sql.NullInt64, userName sql.NullString, avatar sql.NullString) (int64, error) {
	return 1, nil
}

// Update returns no error.
func (d *DummyUser) Update(id int64, email string, provider sql.NullString, token sql.NullString, uuid sql.NullInt64, userName sql.NullString, avatar sql.NullString) error {
	return nil
}

// UpdatePassword returns no error.
func (d *DummyUser) UpdatePassword(id int64, pasword string) error {
	return nil
}
