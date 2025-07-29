package data

// TODO: Add methods to the repository interface
type Repository interface {
	InsertPatrol(*Patrol) (string, error)
	GetAllPatrols(*GetPatrolInfoRequest) ([]*Patrol, error)
	PutPatrolInfo(*UpdatePatrolInfoRequest) error
	PatchPatrolInfo(*UpdatePatrolInfoRequest) error
}
