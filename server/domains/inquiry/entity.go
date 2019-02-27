package inquiry

// Inquiry is a entity for inquiry.
type Inquiry struct {
	ID             int64
	Email          string
	Name           string
	Message        string
	infrastructure Repository
}

// Repository defines repository interface.
type Repository interface {
	Create(string, string, string) (int64, error)
}

// New returns a inquiry struct.
func New(id int64, email, name, message string, infrastructure Repository) *Inquiry {
	i := &Inquiry{
		id,
		email,
		name,
		message,
		infrastructure,
	}
	return i
}

// Create a inquiry entity, and returns the latest inquiry object.
func (i *Inquiry) Create() error {
	id, err := i.infrastructure.Create(i.Email, i.Name, i.Message)
	if err != nil {
		return err
	}
	i.ID = id
	return nil
}
