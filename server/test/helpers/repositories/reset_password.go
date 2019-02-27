package repositories

import "time"

type DummyResetPassword struct{}

func (r *DummyResetPassword) Authenticate(id int64, token string) error {
	return nil
}

func (r *DummyResetPassword) FindAvailable(id int64, token string) (int64, int64, string, time.Time, error) {
	return id, 1, token, time.Now(), nil
}

func (r *DummyResetPassword) Find(id int64) (int64, int64, string, time.Time, error) {
	return id, 1, "", time.Now(), nil
}

func (r *DummyResetPassword) Create(id int64, token string, expiresAt time.Time) (int64, error) {
	return 1, nil
}

func (r *DummyResetPassword) UpdateExpire(id int64) error {
	return nil
}
