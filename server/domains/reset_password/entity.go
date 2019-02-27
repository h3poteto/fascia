package reset_password

import (
	"time"

	"github.com/h3poteto/fascia/server/domains/entities/user"
)

// ResetPassword has a reset password model object
type ResetPassword struct {
	ID             int64
	UserID         int64
	Token          string
	ExpiresAt      time.Time
	infrastructure Repository
}

// Repository defines infrastructure.
type Repository interface {
	Authenticate(int64, string) error
	FindAvailable(int64, string) (int64, int64, string, time.Time, error)
	Find(int64) (int64, int64, string, time.Time, error)
	Create(int64, string, time.Time) (int64, error)
	UpdateExpire(int64) error
}

// New returns a reset password entity
func New(id, userID int64, token string, expiresAt time.Time, infrastructure Repository) *ResetPassword {
	r := &ResetPassword{
		id,
		userID,
		token,
		expiresAt,
		infrastructure,
	}
	return r
}

// Create creates a repository record.
func (r *ResetPassword) Create() error {
	id, err := r.infrastructure.Create(r.UserID, r.Token, r.ExpiresAt)
	if err != nil {
		return err
	}
	r.ID = id
	return nil
}

// User returns a owner user entity
func (r *ResetPassword) User() (*user.User, error) {
	u, err := user.Find(r.UserID)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// UpdateExpire change expire to now.
func (r *ResetPassword) UpdateExpire() error {
	return r.infrastructure.UpdateExpire(r.ID)
}
