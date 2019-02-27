package repositories

import "time"

// DummyResetPassword is a struct for test repository.
type DummyResetPassword struct{}

// Authenticate returns no error for tests.
func (r *DummyResetPassword) Authenticate(id int64, token string) error {
	return nil
}

// FindAvailable returns dummy date for tests.
func (r *DummyResetPassword) FindAvailable(id int64, token string) (int64, int64, string, time.Time, error) {
	return id, 1, token, time.Now(), nil
}

// Find returns dummy data for tests.
func (r *DummyResetPassword) Find(id int64) (int64, int64, string, time.Time, error) {
	return id, 1, "", time.Now(), nil
}

// Create returns no error.
func (r *DummyResetPassword) Create(id int64, token string, expiresAt time.Time) (int64, error) {
	return 1, nil
}

// UpdateExpire returns no error.
func (r *DummyResetPassword) UpdateExpire(id int64) error {
	return nil
}
