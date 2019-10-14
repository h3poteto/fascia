package inquiry

// Repository defines repository interface.
type Repository interface {
	Create(string, string, string) (int64, error)
}
