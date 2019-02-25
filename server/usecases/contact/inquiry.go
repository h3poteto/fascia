package contact

import (
	"database/sql"

	"github.com/h3poteto/fascia/lib/modules/database"
	entity "github.com/h3poteto/fascia/server/domains/entities/inquiry"
	repo "github.com/h3poteto/fascia/server/infrastructures/inquiry"
)

func InjectDB() *sql.DB {
	return database.SharedInstance().Connection
}

func InjectInquiryRepository() entity.InquiryRepository {
	return repo.New(InjectDB())
}

func CreateInquiry(email, name, message string) (*entity.Inquiry, error) {
	inquiry := entity.New(0, email, name, message, InjectInquiryRepository())
	err := inquiry.Create()
	if err != nil {
		return nil, err
	}
	return inquiry, nil
}
