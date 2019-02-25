package contact

import (
	"database/sql"

	"github.com/h3poteto/fascia/lib/modules/database"
	domain "github.com/h3poteto/fascia/server/domains/inquiry"
	repo "github.com/h3poteto/fascia/server/infrastructures/inquiry"
)

// InjectDB set DB connection from connection pool.
func InjectDB() *sql.DB {
	return database.SharedInstance().Connection
}

// InjectInquiryRepository inject db connection and return repository instance.
func InjectInquiryRepository() domain.Repository {
	return repo.New(InjectDB())
}

// CreateInquiry create a new inquiry.
func CreateInquiry(email, name, message string) (*domain.Inquiry, error) {
	inquiry := domain.New(0, email, name, message, InjectInquiryRepository())
	err := inquiry.Create()
	if err != nil {
		return nil, err
	}
	return inquiry, nil
}
