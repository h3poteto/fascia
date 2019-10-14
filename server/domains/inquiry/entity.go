package inquiry

// Inquiry is a entity for inquiry.
type Inquiry struct {
	ID      int64
	Email   string
	Name    string
	Message string
}

// New returns a inquiry struct.
func New(id int64, email, name, message string) *Inquiry {
	i := &Inquiry{
		id,
		email,
		name,
		message,
	}
	return i
}
