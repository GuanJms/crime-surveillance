package data

type Repository interface {
	GetAllCrimes() ([]*Crime, error)
	InsertNewCrime(*Crime) error
}
