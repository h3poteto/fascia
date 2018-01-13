package user

import (
	"github.com/h3poteto/fascia/server/entities/project"
	"github.com/h3poteto/fascia/server/infrastructures/user"
)

// Find returns a user entity
func Find(id int64) (*User, error) {
	u := &User{
		ID: id,
	}
	if err := u.reload(); err != nil {
		return nil, err
	}
	return u, nil
}

// FindByEmail returns a user entity
func FindByEmail(email string) (*User, error) {
	infrastructure, err := user.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	u := &User{
		infrastructure: infrastructure,
	}
	if err := u.reload(); err != nil {
		return nil, err
	}
	return u, nil
}

// Projects list up projects related a user
func (u *User) Projects() ([]*project.Project, error) {
	return project.Projects(u.ID)
}
