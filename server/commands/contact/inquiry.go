package contact

import (
	"github.com/h3poteto/fascia/server/entities/inquiry"
)

// CreateInquiry has inquiry entity.
type CreateInquiry struct {
	InquiryEntity *inquiry.Inquiry
}

// InitCreateInquiry initialize a CreateInquiry struct.
func InitCreateInquiry(id int64, email, name, message string) *CreateInquiry {
	return &CreateInquiry{
		InquiryEntity: inquiry.New(id, email, name, message),
	}
}

// Run save a inquiry entity.
func (i *CreateInquiry) Run() error {
	entity, err := i.InquiryEntity.Save()
	if err != nil {
		return err
	}
	i.InquiryEntity = entity
	return nil
}
