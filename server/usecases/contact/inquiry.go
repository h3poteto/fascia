package contact

import (
	"database/sql"

	"github.com/h3poteto/fascia/lib/modules/database"
	entity "github.com/h3poteto/fascia/server/domains/entities/inquiry"
	repo "github.com/h3poteto/fascia/server/infrastructures/inquiry"
)

// InjectDB set DB connection from connection pool.
func InjectDB() *sql.DB {
	return database.SharedInstance().Connection
}

// InjectInquiryRepository inject db connection and return repository instance.
func InjectInquiryRepository() entity.Repository {
	return repo.New(InjectDB())
}

// CreateInquiry create a new inquiry.
func CreateInquiry(email, name, message string) (*entity.Inquiry, error) {
	inquiry := entity.New(0, email, name, message, InjectInquiryRepository())
	err := inquiry.Create()
	if err != nil {
		return nil, err
	}
	return inquiry, nil
}
