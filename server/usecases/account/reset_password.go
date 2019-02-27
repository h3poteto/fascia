package account

import (
	"database/sql"

	"github.com/h3poteto/fascia/lib/modules/database"
	user "github.com/h3poteto/fascia/server/domains/entities/user"
	domain "github.com/h3poteto/fascia/server/domains/reset_password"
	repo "github.com/h3poteto/fascia/server/infrastructures/reset_password"
)

// InjectDB set DB connection from connection pool.
func InjectDB() *sql.DB {
	return database.SharedInstance().Connection
}

// InjectResetPasswordRepository inject db connection and return repository instance.
func InjectResetPasswordRepository() domain.Repository {
	return repo.New(InjectDB())
}

// ChangeUserPassword change password in user, and expire reset password.
func ChangeUserPassword(id int64, token string, password string) (*user.User, error) {
	reset, err := domain.FindAvailable(id, token, InjectResetPasswordRepository())
	if err != nil {
		return nil, err
	}
	u, err := reset.User()
	if err != nil {
		return nil, err
	}
	if err := u.UpdatePassword(password, nil); err != nil {
		return nil, err
	}
	if err := reset.UpdateExpire(); err != nil {
		return nil, err
	}
	return reset.User()
}

// GenerateResetPassword generates a token and create a new reset password entity.
func GenerateResetPassword(userID int64, email string) (*domain.ResetPassword, error) {
	return domain.GenerateResetPassword(userID, email, InjectResetPasswordRepository())
}

// AuthenticateResetPassword authenticate a reset password.
func AuthenticateResetPassword(id int64, token string) error {
	return domain.Authenticate(id, token, InjectResetPasswordRepository())
}
