package reset_password

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// FindAvailable find available reset password entity
func FindAvailable(targetID int64, targetToken string, infrastructure Repository) (*ResetPassword, error) {
	id, userID, token, expiresAt, err := infrastructure.FindAvailable(targetID, targetToken)
	if err != nil {
		return nil, err
	}
	r := New(id, userID, token, expiresAt, infrastructure)
	return r, nil
}

// Find find a reset password entity by id.
func Find(targetID int64, infrastructure Repository) (*ResetPassword, error) {
	id, userID, token, expiresAt, err := infrastructure.Find(targetID)
	if err != nil {
		return nil, err
	}
	r := New(id, userID, token, expiresAt, infrastructure)
	return r, nil
}

// GenerateResetPassword generate new token and return a new reset password entity, and save it.
func GenerateResetPassword(userID int64, email string, infrastructure Repository) (*ResetPassword, error) {
	// Generate token using md5.
	h := md5.New()
	_, err := io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		return nil, errors.Wrap(err, "token generate error")
	}
	_, err = io.WriteString(h, email)
	if err != nil {
		return nil, errors.Wrap(err, "token generate error")
	}
	token := fmt.Sprintf("%x", h.Sum(nil))

	reset := New(0, userID, token, time.Now().AddDate(0, 0, 1), infrastructure)
	if err := reset.Create(); err != nil {
		return nil, err
	}
	return reset, nil
}

// Authenticate authenticate the reset password record.
func Authenticate(id int64, token string, infrastructure Repository) error {
	return infrastructure.Authenticate(id, token)
}
