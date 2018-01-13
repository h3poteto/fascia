package reset_password

import (
	"github.com/h3poteto/fascia/server/entities/user"
	"github.com/h3poteto/fascia/server/infrastructures/reset_password"
)

// FindAvailable find available reset password entity
func FindAvailable(id int64, token string) (*ResetPassword, error) {
	i, err := reset_password.FindAvailable(id, token)
	if err != nil {
		return nil, err
	}
	r := &ResetPassword{
		infrastructure: i,
	}
	if err := r.reload(); err != nil {
		return nil, err
	}
	return r, nil
}

// Find find a reset password entity by id.
func Find(id int64) (*ResetPassword, error) {
	i, err := reset_password.Find(id)
	if err != nil {
		return nil, err
	}
	r := &ResetPassword{
		infrastructure: i,
	}
	if err := r.reload(); err != nil {
		return nil, err
	}
	return r, nil
}

// User returns a owner user entity
func (r *ResetPassword) User() (*user.User, error) {
	u, err := user.Find(r.UserID)
	if err != nil {
		return nil, err
	}
	return u, nil
}
