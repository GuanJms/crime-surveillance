package data

type Repository interface {
	GetAllCrimes() ([]*Crime, error)
	InsertNewCrime(*Crime) error
	PutCrime(*Crime) error         // Update the whole exisitng crime
	PatchCrime(*CrimeUpdate) error // Partial udpate
	DeleteCrime(id string) error   // delete by the crime id
}
