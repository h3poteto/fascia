package contact

import (
	"github.com/h3poteto/fascia/server/entities/inquiry"
)

type CreateInquiry struct {
	InquiryEntity *inquiry.Inquiry
}

func InitCreateInquiry(id int64, email, name, message string) *CreateInquiry {
	return &CreateInquiry{
		InquiryEntity: inquiry.New(id, email, name, message),
	}
}

func (i *CreateInquiry) Run() error {
	entity, err := i.InquiryEntity.Save()
	if err != nil {
		return err
	}
	i.InquiryEntity = entity
	return nil
}
