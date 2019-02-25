package inquiry

// Inquiry is a entity for inquiry.
type Inquiry struct {
	ID             int64
	Email          string
	Name           string
	Message        string
	infrastructure InquiryRepository
}

type InquiryRepository interface {
	Create(string, string, string) (int64, error)
}

// New returns a inquiry struct.
func New(id int64, email, name, message string, infrastructure InquiryRepository) *Inquiry {
	i := &Inquiry{
		id,
		email,
		name,
		message,
		infrastructure,
	}
	return i
}

// Save a inquiry entity, and returns the latest inquiry object.
func (i *Inquiry) Create() error {
	id, err := i.infrastructure.Create(i.Email, i.Name, i.Message)
	if err != nil {
		return err
	}
	i.ID = id
	return nil
}
