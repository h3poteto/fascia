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
	repo := InjectInquiryRepository()
	id, err := repo.Create(email, name, message)
	if err != nil {
		return nil, err
	}
	inquiry := domain.New(id, email, name, message)
	return inquiry, nil
}
